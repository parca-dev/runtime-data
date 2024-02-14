package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	fSet := flag.NewFlagSet("structlayout", flag.ExitOnError)

	var outputDir string
	fSet.StringVar(&outputDir, "output", "", "output directory to write the merged layout file")
	fSet.StringVar(&outputDir, "o", "", "output directory to write the merged layout file (shorthand)")

	fSet.Usage = func() {
		fmt.Printf("usage: mergelayout -o outputDir <path-to-layout-files>\n")
		fmt.Printf("e.g: mergelayout -o /tmp/merged '/tmp/python/*.yaml'\n\n")
		fmt.Println("flags:")
		fSet.PrintDefaults()
	}

	if err := fSet.Parse(os.Args[1:]); err != nil {
		logger.Error("failed to parse flags", "err", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fSet.Usage()
		os.Exit(1)
	}

	var inputs []string
	if len(fSet.Args()) == 0 {
		logger.Error("at least one input file is required")
		os.Exit(1)
	}
	inputs = fSet.Args()

	if len(inputs) == 1 {
		logger.Info("single input specified, using glob pattern", "pattern", inputs[0])
		input := inputs[0]
		matches, err := filepath.Glob(input)
		if err != nil {
			logger.Error("failed to glob input", "err", err)
			os.Exit(1)
		}

		if len(matches) == 0 {
			logger.Error("no files found", "pattern", input)
			os.Exit(1)
		}

		inputs = matches
		logger.Info("found files", "files", matches, "pattern", input)
	} else {
		logger.Info("files are specified  as input", "files", inputs)
	}

	if outputDir == "" {
		outputDir = "."
	}
	if err := mergeLayoutFiles(logger, inputs, outputDir); err != nil {
		logger.Error("failed to merge files", "err", err)
		os.Exit(1)
	}

	logger.Info("done", "output directory", outputDir)
}

func mergeLayoutFiles(logger *slog.Logger, inputFiles []string, output string) error {
	// Read all the input files and store them in a map with the version as the key.
	versionedLayouts := map[runtimedata.Version]*runtimedata.DataWithVersion{}
	for _, file := range inputFiles {
		logger.Info("reading file", "file", file)
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var withVersion runtimedata.DataWithVersion
		if err := yaml.Unmarshal(data, &withVersion); err != nil {
			return err
		}

		versionedLayouts[withVersion.Version] = &withVersion
	}

	// Versions need to be sorted in ascending order.
	versionKeys := maps.Keys(versionedLayouts)
	sort.Slice(versionKeys, func(i, j int) bool {
		return convertVersion(versionKeys[i]).LessThan(convertVersion(versionKeys[j]))
	})

	var (
		minVersion      *semver.Version
		maxVersion      *semver.Version
		currentLayout   map[string]any
		outputData      = map[string]any{}
		addVersionRange = func() {
			var rawConstr string
			if maxVersion.Equal(minVersion) {
				rawConstr = fmt.Sprintf("= %s", minVersion)
			} else {
				rawConstr = fmt.Sprintf("%s - %s", minVersion, maxVersion)
			}
			_, err := semver.NewConstraint(rawConstr)
			if err != nil {
				panic(err)
			}
			// constr.String() is not used here because it doesn't provide valid file names.
			outputData[rawConstr] = currentLayout
		}
	)
	for _, v := range versionKeys {
		currentVersion := convertVersion(v)
		data := versionedLayouts[v].Data
		if minVersion == nil {
			minVersion = currentVersion
			maxVersion = currentVersion
			currentLayout = data
			continue
		}

		if reflect.DeepEqual(currentLayout, data) {
			maxVersion = currentVersion
			continue
		}

		addVersionRange()
		minVersion = currentVersion
		maxVersion = currentVersion
		currentLayout = data
	}

	// Add the last version range if it's not there already.
	// If there's only one version, it's not added in the loop.
	addVersionRange()

	for versionRange, data := range outputData {
		outputFilePath := filepath.Join(output, fmt.Sprintf("%s.yaml", versionRange))
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}

		encoder := yaml.NewEncoder(outputFile)
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("failed to encode layout: %w", err)
		}
		if err := encoder.Close(); err != nil {
			return fmt.Errorf("failed to encode layout: %w", err)
		}
	}

	return nil
}

func convertVersion(v runtimedata.Version) *semver.Version {
	return semver.New(v.Major, v.Minor, v.Patch, "", "")
}
