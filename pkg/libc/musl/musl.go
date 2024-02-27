package musl

import (
	"embed"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"

	"github.com/parca-dev/runtime-data/pkg/libc"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
)

const layoutDir = "layout"

var (
	//go:embed layout/*/*.yaml
	generatedLayouts embed.FS
	structLayouts    = map[runtimedata.Key]*libc.Layout{}
	once             = &sync.Once{}
)

func init() {
	var err error
	structLayouts, err = loadLayouts()
	if err != nil {
		panic(err)
	}
}

func loadLayouts() (map[runtimedata.Key]*libc.Layout, error) {
	var err error
	once.Do(func() {
		entries, err := generatedLayouts.ReadDir(filepath.Join(layoutDir, runtime.GOARCH))
		if err != nil {
			return
		}
		var i int
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			var data []byte
			data, err = generatedLayouts.ReadFile(filepath.Join(layoutDir, runtime.GOARCH, entry.Name()))
			if err != nil {
				return
			}
			ext := filepath.Ext(entry.Name())
			// Filter out non-yaml files.
			if ext != ".yaml" && ext != ".yml" {
				continue
			}
			var lyt libc.Layout
			if err = yaml.Unmarshal(data, &lyt); err != nil {
				return
			}
			rawConstraint := strings.TrimSuffix(entry.Name(), ext)
			constr, err := semver.NewConstraint(rawConstraint)
			if err != nil {
				return
			}
			key := runtimedata.Key{Index: i, Constraint: constr.String()}
			structLayouts[key] = &lyt
			i++
		}
	})
	return structLayouts, err
}

func getLayoutForArch(v *semver.Version, arch string) (runtimedata.Key, *libc.Layout, error) {
	entries, err := generatedLayouts.ReadDir(filepath.Join(layoutDir, arch))
	if err != nil {
		return runtimedata.Key{}, nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		var data []byte
		data, err = generatedLayouts.ReadFile(filepath.Join(layoutDir, arch, entry.Name()))
		if err != nil {
			return runtimedata.Key{}, nil, err
		}
		ext := filepath.Ext(entry.Name())
		// Filter out non-yaml files.
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		var lyt libc.Layout
		if err = yaml.Unmarshal(data, &lyt); err != nil {
			return runtimedata.Key{}, nil, err
		}
		rawConstraint := strings.TrimSuffix(entry.Name(), ext)
		constr, err := semver.NewConstraint(rawConstraint)
		if err != nil {
			return runtimedata.Key{}, nil, err
		}
		key := runtimedata.Key{Constraint: constr.String()}
		if constr.Check(v) {
			return key, &lyt, nil
		}
	}
	return runtimedata.Key{}, nil, errors.New("not found")
}

// GetLayout returns the layout for the given version.
func GetLayout(v *semver.Version) (runtimedata.Key, *libc.Layout, error) {
	for k, l := range structLayouts {
		constr, err := semver.NewConstraint(k.Constraint)
		if err != nil {
			return k, nil, err
		}
		if constr.Check(v) {
			return k, l, nil
		}
	}
	return runtimedata.Key{}, nil, errors.New("not found")
}

// GetLayouts returns all the layouts.
func GetLayouts() (map[runtimedata.Key]*libc.Layout, error) {
	return structLayouts, nil
}
