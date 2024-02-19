package main

import (
	"debug/dwarf"
	"debug/elf"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
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
		fmt.Printf("e.g: structlayout -r python -v 3.9.5 /usr/bin/python3.9\n\n")
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
		layoutMap       runtimedata.LayoutMap
		initialStateMap any
		outputDir       = givenOutputDir
	)
	switch runtime {
	case "python":
		layoutMap = python.DataMapForLayout(version)
		initialStateMap = python.DataMapForInitialState(version)
		if outputDir == "" {
			// Base output directory for python is pkg/python.
			outputDir = "pkg/python"
		}
	case "ruby":
		layoutMap = ruby.DataMapForLayout(version)
		if outputDir == "" {
			outputDir = "pkg/ruby"
		}
	case "libc":
		// TODO(kakkoyun): Change depending on the libc implementation. e.g musl, glibc, etc.
		// layoutMap = libc.DataMapForLayout(version)
		// if outputDir == "" {
		// 	outputDir = "pkg/libc/layout"
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
		input = fSet.Arg(0)
	)
	dwarfData, err := dwarfDataFromELF(input)
	if err != nil {
		logger.Error("failed to read DWARF data", "err", err)
		os.Exit(1)
	}

	output := filepath.Join(outputDir, "layout", fmt.Sprintf("%s_%s.yaml", runtime, sanitizeIdentifier(version)))
	if err := processAndWriteLayout(dwarfData, output, version, layoutMap); err != nil {
		logger.Error("failed to write layout", "err", err)
		os.Exit(1)
	}
	logger.Info("layout file written", "file", output)

	if isNil(initialStateMap) {
		logger.Info("no initial state map found, skipping initial state generation")
		os.Exit(0)
	}

	output = filepath.Join(outputDir, "initialstate", fmt.Sprintf("%s_%s.yaml", runtime, sanitizeIdentifier(version)))
	if err := processAndWriteInitialState(dwarfData, output, version, initialStateMap); err != nil {
		logger.Error("failed to write initial state", "err", err)
		os.Exit(1)
	}
	logger.Info("initial state file written", "file", output)
}

// processAndWriteLayout processes the given ELF file and writes the layout to the given output file.
func processAndWriteLayout(dwarfData *dwarf.Data, output string, version string, layoutMap runtimedata.LayoutMap) error {
	dm, err := datamap.New(layoutMap)
	if err != nil {
		return fmt.Errorf("failed to create data map: %w", err)
	}

	if err := dm.ReadFromDWARF(dwarfData); err != nil {
		return fmt.Errorf("failed to extract struct layout from DWARF data: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	// Extremely in-efficient and hacky but it should work for now.
	withVersion, err := runtimedata.WithVersion(version, convertToMapOfAny(layoutMap.Layout()))
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

// processAndWriteInitialState processes the given ELF file and writes the initial state to the given output file.
func processAndWriteInitialState(dwarfData *dwarf.Data, output string, version string, initialStateMap any) error {
	dm, err := datamap.New(initialStateMap)
	if err != nil {
		return fmt.Errorf("failed to create data map: %w", err)
	}

	if err := dm.ReadFromDWARF(dwarfData); err != nil {
		return fmt.Errorf("failed to extract struct layout from DWARF data: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	// Extremely in-efficient and hacky but it should work for now.
	withVersion, err := runtimedata.WithVersion(version, convertToMapOfAny(initialStateMap))
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

// dwarfDataFromELF returns the DWARF data from the given ELF file.
func dwarfDataFromELF(input string) (*dwarf.Data, error) {
	ef, err := elf.Open(input)
	if err != nil {
		return nil, fmt.Errorf("failed to open ELF file: %w", err)
	}
	defer ef.Close()

	dwarfData, err := ef.DWARF()
	if err != nil {
		return nil, fmt.Errorf("failed to read DWARF info: %w", err)
	}

	return dwarfData, nil
}

// sanitizeIdentifier sanitizes the identifier to be used as a filename.
func sanitizeIdentifier(identifier string) string {
	return strings.TrimPrefix(strings.ReplaceAll(identifier, ".", "_"), "v")
}

// convertToMapOfAny converts the given struct to a map of string to any.
func convertToMapOfAny(v interface{}) map[string]any {
	// Marshal and unmarshal to convert the struct to map[string]any.
	blob, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}

	var anyMap map[string]any
	if err := yaml.Unmarshal(blob, &anyMap); err != nil {
		panic(err)
	}

	return anyMap
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return val.IsNil()
	}
	return false
}
