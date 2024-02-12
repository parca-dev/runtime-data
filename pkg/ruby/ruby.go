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
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

type Key struct {
	index      int
	constraint string
}

const versionDir = "versions"

var (
	//go:embed versions/*.yaml
	generatedLayouts embed.FS
	versions         = map[Key]*Layout{}
	once             = &sync.Once{}
)

func init() {
	var err error
	versions, err = loadVersionLayouts()
	if err != nil {
		panic(err)
	}
}

func loadVersionLayouts() (map[Key]*Layout, error) {
	var err error
	once.Do(func() {
		entries, err := generatedLayouts.ReadDir(versionDir)
		if err != nil {
			return
		}
		var i int
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			var data []byte
			data, err = generatedLayouts.ReadFile(filepath.Join(versionDir, entry.Name()))
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
			key := Key{index: i, constraint: constr.String()}
			versions[key] = &lyt
			i++
		}
	})
	return versions, err
}

// GetLayout returns the matching layout for the given version.
func GetLayout(v *semver.Version) (Key, *Layout, error) {
	versions, err := loadVersionLayouts()
	if err != nil {
		return Key{}, nil, err
	}
	for k, l := range versions {
		constr, err := semver.NewConstraint(k.constraint)
		if err != nil {
			return k, nil, err
		}
		if constr.Check(v) {
			return k, l, nil
		}
	}
	return Key{}, nil, errors.New("not found")
}

// GetLayouts returns all the layouts for the supported versions.
func GetLayouts() (map[Key]*Layout, error) {
	versions, err := loadVersionLayouts()
	if err != nil {
		return nil, err
	}
	return versions, nil
}
