package main

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/blakesmith/ar"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/ulikunitz/xz"
	"golang.org/x/net/html"
)

const (
	// DefaultUbuntuBaseURL is the default base URL to download deb files from.
	DefaultUbuntuBaseURL = "http://archive.ubuntu.com/ubuntu/pool/main"
	// DefaultDebianBaseURL is the default base URL to download deb files from.
	DefaultDebianBaseURL = "http://ftp.debian.org/debian/pool/main"

	// DefaultBaseURL is the default base URL to download deb files from.
	DefaultBaseURL = DefaultDebianBaseURL
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	fSet := ff.NewFlagSet("debdownload")
	var (
		outputDir         = fSet.String('o', "output", "", "output directory to write the downloaded deb files")
		tempDir           = fSet.String('t', "temp-dir", "", "temporary directory to download deb files")
		url               = fSet.String('u', "url", "", "URL to download deb files from")
		pkgName           = fSet.String('p', "package", "", "package name to download")
		architectures     = fSet.StringList('a', "arch", "architectures to download")
		versionConstraint = fSet.String('c', "constraint", "", "version constraints to download")
	)
	if err := ff.Parse(fSet, os.Args[1:]); err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fSet))
		if !errors.Is(err, ff.ErrHelp) {
			fmt.Printf("err=%v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *outputDir == "" {
		logger.Error("output directory is required")
		os.Exit(1)
	}

	if *tempDir == "" {
		*tempDir = os.TempDir()
	}

	if *url == "" {
		*url = DefaultBaseURL
	}

	if *pkgName == "" {
		logger.Error("package name is required")
		os.Exit(1)
	}

	if len(*architectures) == 0 {
		*architectures = []string{"amd64", "arm64"}
	}

	cli := &cli{logger: logger}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	packages, err := cli.list(ctx, *url, *pkgName, *architectures, *versionConstraint)
	if err != nil {
		logger.Error("failed to list packages", "err", err)
		os.Exit(1)
	}

	interimDir := filepath.Join(*tempDir, *pkgName)
	if err := os.MkdirAll(interimDir, 0755); err != nil {
		logger.Error("failed to create temp directory", "err", err)
		os.Exit(1)
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].version.GreaterThan(packages[j].version)
	})

	if err := cli.download(ctx, packages, interimDir); err != nil {
		logger.Error("failed to download packages", "err", err)
		os.Exit(1)
	}

	targetDir := filepath.Join(*outputDir, *pkgName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		logger.Error("failed to create output directory", "err", err)
		os.Exit(1)
	}

	if err := cli.extract(ctx, packages, targetDir); err != nil {
		logger.Error("failed to extract packages", "err", err)
		os.Exit(1)
	}

	logger.Info("downloaded packages", "outputDir", *outputDir)
}

type cli struct {
	logger *slog.Logger
}

type pkg struct {
	link    string
	name    string
	variant string
	version *semver.Version
	arch    string

	downloadedArchive string
}

func (c *cli) list(ctx context.Context, pkgUrl, pkgName string, architectures []string, versionConstraint string) ([]*pkg, error) {
	c.logger.Info("listing packages", "url", pkgUrl)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pkgUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var matcher = regexp.MustCompile(
		fmt.Sprintf(`(%s)(-.*)?_(.*)_(%s)\.deb`, pkgName, strings.Join(architectures, "|")),
	)
	c.logger.Info("matcher", "pattern", matcher.String())

	var packages []*pkg
	var process func(*html.Node)
	process = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if matches := matcher.FindStringSubmatch(a.Val); len(matches) > 0 {
						version := semver.MustParse(matches[3])
						if versionConstraint != "" {
							match, reasons := mustConstraint(semver.NewConstraint(versionConstraint)).Validate(version)
							if !match {
								c.logger.Info("version does not match", "version", version, "reasons", reasons)
								c.logger.Info("see: https://github.com/Masterminds/semver?tab=readme-ov-file#checking-version-constraints")
								continue
							}
						}
						packages = append(packages, &pkg{
							link:    must(url.JoinPath(pkgUrl, a.Val)),
							name:    matches[1],
							variant: strings.TrimPrefix(matches[2], "-"),
							version: version,
							arch:    matches[4],
						})
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			process(c)
		}
	}
	process(doc)
	return packages, nil
}

func (c *cli) download(ctx context.Context, packages []*pkg, tempDir string) error {
	c.logger.Info("downloading packages", "tempDir", tempDir)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for _, p := range packages {
		shortVersion := fmt.Sprintf("%d.%d.%d", p.version.Major(), p.version.Minor(), p.version.Patch())
		var target string
		if p.variant != "" {
			target = filepath.Join(tempDir, fmt.Sprintf("%s-%s_%s_%s.deb", p.name, p.variant, shortVersion, p.arch))
		} else {
			target = filepath.Join(tempDir, fmt.Sprintf("%s_%s_%s.deb", p.name, shortVersion, p.arch))
		}
		if _, err := os.Stat(target); err == nil {
			c.logger.Info("file already exists", "file", target)
			p.downloadedArchive = target
			continue
		}

		c.logger.Info("downloading package", "link", p.link, "target", target)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.link, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		f, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()

		if _, err := io.Copy(f, resp.Body); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		p.downloadedArchive = target
	}
	return nil
}

func (c *cli) extract(_ context.Context, packages []*pkg, outputDir string) error {
	c.logger.Info("extracting packages", "outputDir", outputDir)

	for _, p := range packages {
		if p.downloadedArchive == "" {
			continue
		}

		if !strings.HasSuffix(p.downloadedArchive, ".deb") {
			continue
		}

		f, err := os.Open(p.downloadedArchive)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		c.logger.Info("extracting package", "file", f.Name())

		var variant string
		if p.variant != "" {
			variant = p.variant
		} else {
			variant = "main"
		}
		shortVersion := fmt.Sprintf("%d.%d.%d", p.version.Major(), p.version.Minor(), p.version.Patch())
		targetDir := filepath.Join(outputDir, p.arch, shortVersion, variant)
		if _, err := os.Stat(targetDir); err == nil {
			c.logger.Info("file already exists", "file", targetDir)
			continue
		}

		r := ar.NewReader(f)
		for {
			hdr, err := r.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read next entry: %w", err)
			}

			if hdr.Name != "data.tar.xz" {
				continue
			}

			c.logger.Info("extracting tar.xz", "file", hdr.Name)

			content, err := io.ReadAll(r)
			if err != nil {
				return fmt.Errorf("failed to read content: %w", err)
			}

			xzr, err := xz.NewReader(bytes.NewReader(content))
			if err != nil {
				return fmt.Errorf("failed to create xz reader: %w", err)
			}

			tr := tar.NewReader(xzr)
			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("failed to read next entry: %w", err)
				}
				if hdr == nil {
					break
				}

				target := filepath.Join(targetDir, hdr.Name)
				if hdr.FileInfo().IsDir() {
					if err := os.MkdirAll(target, 0755); err != nil {
						return fmt.Errorf("failed to create directory: %w", err)
					}
					continue
				}

				f, err := os.Create(target)
				if err != nil {
					return fmt.Errorf("failed to create file: %w", err)
				}
				defer f.Close()

				if _, err := io.Copy(f, tr); err != nil {
					return fmt.Errorf("failed to copy file: %w", err)
				}
			}
		}
	}
	return nil
}

func must(u string, err error) string {
	if err != nil {
		panic(err)
	}
	return u
}

func mustConstraint(c *semver.Constraints, err error) *semver.Constraints {
	if err != nil {
		panic(err)
	}
	return c
}
