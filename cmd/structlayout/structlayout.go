package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/parca-dev/runtime-data/pkg/datamap"
	"github.com/parca-dev/runtime-data/pkg/libc"
	"github.com/parca-dev/runtime-data/pkg/python"
	"github.com/parca-dev/runtime-data/pkg/ruby"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
)

type Mapper interface {
	Layout() runtimedata.RuntimeData
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	fSet := flag.NewFlagSet("structlayout", flag.ExitOnError)
	var (
		runtime string
		version string
	)
	fSet.StringVar(&runtime, "runtime", "", "name of the pre-defined runtime, e.g. python, ruby, libc")
	fSet.StringVar(&runtime, "r", "", "name of the pre-defined runtime, e.g. python, ruby, libc (shorthand)")
	fSet.StringVar(&version, "version", "", "version of the runtime that the layout to generate, e.g. 3.9.5")
	fSet.StringVar(&version, "v", "", "version of the runtime that the layout to generate, e.g. 3.9.5 (shorthand)")

	fSet.Usage = func() {
		fmt.Printf("fsage: structlayout [flags] <path-to-elf>\n")
		fmt.Printf("e.g: structlayout -m python -v 3.9.5 /usr/bin/python3.9\n\n")
		fmt.Println("flags:")
		fSet.PrintDefaults()
	}
	if err := fSet.Parse(os.Args[1:]); err != nil {
		logger.Error("failed to parse flags", "err", err)
		os.Exit(1)
	}
	if len(os.Args) < 4 {
		fSet.Usage()
		os.Exit(1)
	}

	var (
		layoutMap runtimedata.LayoutMap
		outputDir string
	)
	switch runtime {
	case "python":
		layoutMap = python.DataMapForVersion(version)
		outputDir = "pkg/python/versions"
	case "ruby":
		layoutMap = ruby.DataMapForVersion(version)
		outputDir = "pkg/ruby/versions"
	case "libc":
		// TODO(kakkoyun): Change depending on the libc implementation. e.g musl, glibc, etc.
		layoutMap = libc.DataMapForVersion(version)
		outputDir = "pkg/libc"
	default:
		logger.Error("invalid offset map module", "mod", runtime)
		os.Exit(1)
	}
	if layoutMap == nil {
		logger.Error("unknown version", "version", version)
		os.Exit(1)
	}

	dm, err := datamap.New(layoutMap)
	if err != nil {
		logger.Error("failed to generate query", "err", err)
		os.Exit(1)
	}

	input := fSet.Arg(0)
	ef, err := elf.Open(input)
	if err != nil {
		logger.Error("failed to open ELF file", "path", input, "err", err)
		os.Exit(1)
	}
	defer ef.Close()

	dwarfData, err := ef.DWARF()
	if err != nil {
		logger.Error("failed to read DWARF info", "err", err)
		os.Exit(1)
	}

	if err := dm.ReadFromDWARF(dwarfData); err != nil {
		logger.Error("failed to read DWARF data", "err", err)
		os.Exit(1)
	}

	output := filepath.Join(outputDir, fmt.Sprintf("%s_%s.yaml", runtime, sanitizeIdentifier(version)))

	file, err := os.Create(output)
	if err != nil {
		logger.Error("failed to create output file", "path", output, "err", err)
		os.Exit(1)
	}

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(layoutMap.Layout()); err != nil {
		logger.Error("failed to encode layout", "err", err)
		os.Exit(1)
	}
	if err := encoder.Close(); err != nil {
		logger.Error("failed to close encoder", "err", err)
		os.Exit(1)
	}

	logger.Info("offsets written to file", "path", output)
}

// sanitizeIdentifier sanitizes the identifier to be used as a filename.
func sanitizeIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, ".", "_")
}
