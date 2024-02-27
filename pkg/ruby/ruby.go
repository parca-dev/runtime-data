// Copyright 2023 The Parca Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ruby

import (
	"embed"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/parca-dev/runtime-data/pkg/runtimedata"
	"gopkg.in/yaml.v3"
)

const layoutDir = "layout"

var (
	//go:embed layout/*/*.yaml
	generatedLayouts embed.FS
	structLayouts    = map[runtimedata.Key]runtimedata.RuntimeData{}
	once             = &sync.Once{}
)

func init() {
	var err error
	structLayouts, err = loadLayouts()
	if err != nil {
		panic(err)
	}
}

func loadLayouts() (map[runtimedata.Key]runtimedata.RuntimeData, error) {
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
			var lyt Layout
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

// GetLayout returns the matching layout for the given version.
func GetLayout(v *semver.Version) (runtimedata.Key, runtimedata.RuntimeData, error) {
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

// GetLayouts returns all the layouts for the supported versions.
func GetLayouts() (map[runtimedata.Key]runtimedata.RuntimeData, error) {
	return structLayouts, nil
}
