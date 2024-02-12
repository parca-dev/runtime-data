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
		runtime        string
		version        string
		givenOutputDir string
	)
	fSet.StringVar(&runtime, "runtime", "", "name of the pre-defined runtime, e.g. python, ruby, libc")
	fSet.StringVar(&runtime, "r", "", "name of the pre-defined runtime, e.g. python, ruby, libc (shorthand)")
	fSet.StringVar(&version, "version", "", "version of the runtime that the layout to generate, e.g. 3.9.5")
	fSet.StringVar(&version, "v", "", "version of the runtime that the layout to generate, e.g. 3.9.5 (shorthand)")
	fSet.StringVar(&givenOutputDir, "output", "", "output directory to write the layout file")
	fSet.StringVar(&givenOutputDir, "o", "", "output directory to write the layout file (shorthand)")

	fSet.Usage = func() {
		fmt.Printf("usage: structlayout [flags] <path-to-elf>\n")
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
		outputDir = givenOutputDir
	)
	switch runtime {
	case "python":
		layoutMap = python.DataMapForVersion(version)
		if outputDir == "" {
			outputDir = "pkg/python/versions"
		}
	case "ruby":
		layoutMap = ruby.DataMapForVersion(version)
		if outputDir == "" {
			outputDir = "pkg/ruby/versions"
		}
	case "libc":
		// TODO(kakkoyun): Change depending on the libc implementation. e.g musl, glibc, etc.
		// layoutMap = libc.DataMapForVersion(version)
		// if outputDir == "" {
		// 	outputDir = "pkg/libc/versions"
		// }
	default:
		logger.Error("invalid offset map module", "mod", runtime)
		os.Exit(1)
	}

	if layoutMap == nil {
		logger.Error("unknown version", "version", version)
		os.Exit(1)
	}

	var (
		input  = fSet.Arg(0)
		output = filepath.Join(outputDir, fmt.Sprintf("%s_%s.yaml", runtime, sanitizeIdentifier(version)))
	)
	if err := processAndWriteLayout(input, output, version, layoutMap); err != nil {
		logger.Error("failed to write layout", "err", err)
		os.Exit(1)
	}

	logger.Info("layout written to file", "file", output)
}

// processAndWriteLayout processes the given ELF file and writes the layout to the given output file.
func processAndWriteLayout(input, output string, version string, layoutMap runtimedata.LayoutMap) error {
	ef, err := elf.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open ELF file: %w", err)
	}
	defer ef.Close()

	dwarfData, err := ef.DWARF()
	if err != nil {
		return fmt.Errorf("failed to read DWARF info: %w", err)
	}

	dm, err := datamap.New(layoutMap)
	if err != nil {
		return fmt.Errorf("failed to create data map: %w", err)
	}

	if err := dm.ReadFromDWARF(dwarfData); err != nil {
		return fmt.Errorf("failed to extract struct layout from DWARF data: %w", err)
	}

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	withVersion, err := runtimedata.WithVersion(version, layoutMap.Layout())
	if err != nil {
		return fmt.Errorf("failed to wrap layout with version: %w", err)
	}

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(withVersion); err != nil {
		return fmt.Errorf("failed to encode layout: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close encoder: %w", err)
	}

	return nil
}

// sanitizeIdentifier sanitizes the identifier to be used as a filename.
func sanitizeIdentifier(identifier string) string {
	return strings.TrimPrefix(strings.ReplaceAll(identifier, ".", "_"), "v")
}
