package main

import (
	"bytes"
	"debug/elf"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"

	"github.com/parca-dev/runtime-data/pkg/buildid"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	fSet := ff.NewFlagSet("debdownload")
	var (
		debuginfoDir = fSet.String('d', "debuginfo-dir", "", "directory to write the downloaded debuginfo files")
	)
	if err := ff.Parse(fSet, os.Args[1:]); err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fSet))
		if !errors.Is(err, ff.ErrHelp) {
			fmt.Printf("err=%v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(fSet.GetArgs()) != 1 {
		logger.Error("target is required")
		os.Exit(1)
	}

	target := fSet.GetArgs()[0]

	if *debuginfoDir == "" {
		logger.Error("debuginfo directory is required")
		os.Exit(1)
	}

	cli := &cli{logger: logger}

	logger.Info("running", "target", target, "debuginfoDir", *debuginfoDir)

	if err := cli.run(*debuginfoDir, target); err != nil {
		logger.Error("failed to run", "err", err)
		os.Exit(1)
	}
}

type cli struct {
	logger *slog.Logger
}

func (c *cli) run(debuginfoDir string, targetPath string) error {
	f, err := os.Open(targetPath)
	if err != nil {
		return fmt.Errorf("failed to open target: %w", err)
	}

	buildID, err := buildid.FromFile(f)
	if err != nil {
		return fmt.Errorf("failed to get build ID: %w", err)
	}

	c.logger.Debug("looking for the debuginfo", "buildID", buildID)

	dbgFile, err := c.find(debuginfoDir, buildID, f)
	if err != nil {
		return fmt.Errorf("failed to find debuginfo file: %w", err)
	}

	fmt.Println(dbgFile)
	return nil
}

var errSectionNotFound = errors.New("section not found")

func (c cli) find(root string, buildID string, f *os.File) (string, error) {
	if len(buildID) < 2 {
		return "", errors.New("invalid build ID")
	}

	ef, err := elf.NewFile(f)
	if err != nil {
		return "", fmt.Errorf("failed to open ELF file: %w", err)
	}

	// There are two ways of specifying the separate debuginfo file:
	// 1) The executable contains a debug link that specifies the name of the separate debuginfo file.
	//	The separate debug file’s name is usually executable.debug,
	//	where executable is the name of the corresponding executable file without leading directories (e.g., ls.debug for /usr/bin/ls).
	// 2) The executable contains a build ID, a unique bit string that is also present in the corresponding debuginfo file.
	//  (This is supported only on some operating systems, when using the ELF or PE file formats for binary files and the GNU Binutils.)
	//  The debuginfo file’s name is not specified explicitly by the build ID, but can be computed from the build ID, see below.
	//
	// So, for example, suppose you ask Agent to debug /usr/bin/ls, which has a debug link that specifies the file ls.debug,
	//	and a build ID whose value in hex is abcdef1234.
	//	If the list of the global debug directories includes /usr/lib/debug (which is the default),
	//	then Finder will look for the following debug information files, in the indicated order:
	//
	//		- /usr/lib/debug/.build-id/ab/cdef1234.debug
	//		- /usr/bin/ls.debug
	//		- /usr/bin/.debug/ls.debug
	//		- /usr/lib/debug/usr/bin/ls.debug
	//
	// For further information, see: https://sourceware.org/gdb/onlinedocs/gdb/Separate-Debug-Files.html

	// A debug link is a special section of the executable file named .gnu_debuglink. The section must contain:
	//
	// A filename, with any leading directory components removed, followed by a zero byte,
	//  - zero to three bytes of padding, as needed to reach the next four-byte boundary within the section, and
	//  - a four-byte CRC checksum, stored in the same endianness used for the executable file itself.
	// The checksum is computed on the debugging information file’s full contents by the function given below,
	// passing zero as the crc argument.

	base, crc, err := readDebuglink(ef)
	if err != nil {
		if !errors.Is(err, errSectionNotFound) {
			c.logger.Debug("failed to read debug links", "err", err)
		}
	}

	files := c.generatePaths(root, buildID, f.Name(), base)
	if len(files) == 0 {
		return "", errors.New("failed to generate paths")
	}

	c.logger.Debug("generated paths", "paths", files)

	var found string
	for _, file := range files {
		_, err := os.Stat(file)
		if err == nil {
			found = file
			break
		}
		if os.IsNotExist(err) || errors.Is(err, fs.ErrNotExist) {
			continue
		}
	}

	if found == "" {
		return "", os.ErrNotExist
	}

	if strings.Contains(found, ".build-id") || strings.HasSuffix(found, "/debuginfo") || crc <= 0 {
		return found, nil
	}

	match, err := checkSum(found, crc)
	if err != nil {
		return "", fmt.Errorf("failed to check checksum: %w", err)
	}

	if match {
		return found, nil
	}

	return "", os.ErrNotExist
}

func readDebuglink(ef *elf.File) (string, uint32, error) {
	if sec := ef.Section(".gnu_debuglink"); sec != nil {
		d, err := sec.Data()
		if err != nil {
			return "", 0, err
		}
		parts := bytes.Split(d, []byte{0})
		name := string(parts[0])
		sum := parts[len(parts)-1]
		if len(sum) != 4 {
			return "", 0, errors.New("invalid checksum length")
		}
		crc := ef.FileHeader.ByteOrder.Uint32(sum)
		if crc == 0 {
			return "", 0, errors.New("invalid checksum")
		}
		return name, crc, nil
	}
	return "", 0, errSectionNotFound
}

var debugDirs = []string{
	"/usr/lib/debug",
}

func (c *cli) generatePaths(root, buildID, path, filename string) []string {
	const dbgExt = ".debug"
	if len(filename) == 0 {
		filename = filepath.Base(path)
	}
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = dbgExt
	}
	dbgFilePath := filepath.Join(filepath.Dir(path), strings.TrimSuffix(filename, ext)) + ext

	var files []string
	for _, dir := range debugDirs {
		rel, err := filepath.Rel(root, dbgFilePath)
		if err != nil {
			continue
		}
		files = append(files, []string{
			dbgFilePath,
			filepath.Join(filepath.Dir(path), dbgExt, filepath.Base(dbgFilePath)),
			filepath.Join(root, dir, rel),
			filepath.Join(root, dir, ".build-id", buildID[:2], buildID[2:]) + dbgExt,
			filepath.Join(root, dir, buildID, "debuginfo"),
		}...)
	}
	return files
}

// NOTE: we are within the race condition window, but alas.
func checkSum(path string, crc uint32) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	d, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}
	return crc == crc32.ChecksumIEEE(d), nil
}
