package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"golang.org/x/exp/maps"
	"golang.org/x/net/html"
)

const (
	// https://dl-cdn.alpinelinux.org/alpine/MIRRORS.txt
	DefaultBaseURL = "https://dl-cdn.alpinelinux.org/alpine/latest-stable/main/"
	// FetchListTimeout is the timeout to fetch the list of packages.
	FetchListTimeout = 30 * time.Second

	// DownloadSinglePackageTimeout is the timeout to download a single package.
	DownloadSinglePackageTimeout = 90 * time.Second
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	fSet := ff.NewFlagSet("apkdownload")
	var (
		outputDir         = fSet.String('o', "output", "tmp/bin", "output directory to write the downloaded apk files")
		tempDir           = fSet.String('t', "temp-dir", "tmp/apk", "temporary directory to download deb files")
		url               = fSet.String('u', "url", "", "URL to download apk files from")
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

	cli := &cli{
		logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	packages, err := cli.list(ctx, *url, *pkgName, *architectures, *versionConstraint)
	if err != nil {
		logger.Error("failed to list packages", "err", err)
		os.Exit(1)
	}

	interimDir := filepath.Join(*tempDir, *pkgName)
	if err := os.MkdirAll(interimDir, 0o755); err != nil {
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
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		logger.Error("failed to create output directory", "err", err)
		os.Exit(1)
	}

	if err := cli.extract(ctx, packages, targetDir); err != nil {
		logger.Error("failed to extract packages", "err", err)
		os.Exit(1)
	}

	logger.Info("downloaded packages", "outputDir", *outputDir)
}

type pkg struct {
	link    string
	name    string
	variant string
	version *semver.Version
	arch    string

	downloadedArchive string
}

type cli struct {
	logger *slog.Logger
}

var allowedVariants = map[string]struct{}{
	"dbg": {},
	// "dev": {},
	// "fts":     {},
	// "legacy":  {},
	// "libintl": {},
	// "locales": {},
	// "obstack": {},
	// "utils":   {},
}

func (c *cli) list(ctx context.Context, pkgUrl, pkgName string, architectures []string, versionConstraint string) ([]*pkg, error) {
	packages := map[string]*pkg{}
	for _, arch := range architectures {
		pkgUrl, err := url.JoinPath(pkgUrl, convertArch(arch))
		if err != nil {
			return nil, fmt.Errorf("failed to join URL: %w", err)
		}
		c.logger.Info("listing packages", "url", pkgUrl)

		ctx, cancel := context.WithTimeout(ctx, FetchListTimeout)
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
			if resp.StatusCode == http.StatusNotFound {
				c.logger.Info("this version probably does not have this architecture supported", "arch", arch, "url", pkgUrl)
				continue
			}
			return nil, fmt.Errorf("unexpected status code (%s): %d", pkgUrl, resp.StatusCode)
		}

		doc, err := html.Parse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML: %w", err)
		}

		matcher := regexp.MustCompile(
			fmt.Sprintf(`%s(-.*)?-([0-9]\.[0-9]\.[0-9][0-9]?).*\.apk`, pkgName),
		)
		c.logger.Info("matcher", "pattern", matcher.String())

		key := func(p *pkg) string {
			return strings.Join([]string{p.name, p.variant, shortVersion(p.version), p.arch}, "-")
		}

		var process func(*html.Node)
		process = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						if matches := matcher.FindStringSubmatch(a.Val); len(matches) > 0 {
							version := semver.MustParse(strings.ReplaceAll(matches[2], "~", "+"))
							if versionConstraint != "" {
								match, reasons := mustConstraint(semver.NewConstraint(versionConstraint)).Validate(version)
								if !match {
									c.logger.Info("version does not match", "version", version, "reasons", reasons)
									c.logger.Info("see: https://github.com/Masterminds/semver?tab=readme-ov-file#checking-version-constraints")
									continue
								}
							}
							variant := strings.TrimPrefix(matches[1], "-")
							if variant != "" {
								if _, ok := allowedVariants[variant]; !ok {
									continue
								}
							}
							p := &pkg{
								link:    must(url.JoinPath(pkgUrl, a.Val)),
								name:    pkgName,
								variant: variant,
								version: version,
								arch:    arch,
							}
							if oldPkg, ok := packages[key(p)]; ok {
								if p.version.GreaterThan(oldPkg.version) {
									c.logger.Info("found newer version", "old", oldPkg.version, "new", p.version)
									packages[key(p)] = p
								}
								continue
							}

							packages[key(p)] = p
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				process(c)
			}
		}
		process(doc)
	}
	return maps.Values(packages), nil
}

func shortVersion(v *semver.Version) string {
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

func (c *cli) download(ctx context.Context, packages []*pkg, tempDir string) error {
	c.logger.Info("downloading packages", "tempDir", tempDir)

	for _, p := range packages {
		var (
			target  string
			version = p.version.String()
		)
		if p.variant != "" {
			target = filepath.Join(tempDir, fmt.Sprintf("%s-%s_%s_%s.apk", p.name, p.variant, version, p.arch))
		} else {
			target = filepath.Join(tempDir, fmt.Sprintf("%s_%s_%s.apk", p.name, version, p.arch))
		}
		if _, err := os.Stat(target); err == nil {
			c.logger.Info("file already exists", "file", target)
			p.downloadedArchive = target
			continue
		}

		c.logger.Info("downloading package", "link", p.link, "target", target)

		ctx, cancel := context.WithTimeout(ctx, DownloadSinglePackageTimeout)
		defer cancel()

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

func (c *cli) extract(ctx context.Context, packages []*pkg, outputDir string) error {
	c.logger.Info("extracting packages", "outputDir", outputDir)

	for _, p := range packages {
		if p.downloadedArchive == "" {
			continue
		}

		if !strings.HasSuffix(p.downloadedArchive, ".apk") {
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

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(f); err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		b := buf.Bytes()
		var offsets []int
		for i := range b {
			if readGZIPHeader(b[i:]) {
				offsets = append(offsets, i)
			}
		}
		block := b[offsets[2]:]
		br := bytes.NewReader(block)
		r, err := gzip.NewReader(br)
		if err != nil {
			return err
		}
		defer r.Close()

		tr := tar.NewReader(r)
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
				if err := os.MkdirAll(target, 0o755); err != nil {
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
	return nil
}

// Signature bytes for finding GZIP header
const (
	GzipID1     = 0x1f
	GzipID2     = 0x8b
	GzipDeflate = 8
)

// readGZIPHeader reads the header of a gzip file if found.
func readGZIPHeader(buf []byte) bool {
	if buf[0] != GzipID1 || buf[1] != GzipID2 || buf[2] != GzipDeflate {
		return false
	}
	return true
}

func convertArch(arch string) string {
	switch arch {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	}
	return arch
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
