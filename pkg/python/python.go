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
package python

import (
	"embed"
	"errors"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

const (
	layouts      = "layout"
	initialState = "initialstate"
)

type Key struct {
	index      int
	constraint string
}

var (
	//go:embed layout/*.yaml
	generatedLayouts embed.FS
	//go:embed initialstate/*
	generatedState embed.FS

	structLayouts = map[Key]*Layout{}
	once          = &sync.Once{}
)

func init() {
	var err error
	structLayouts, err = loadLayouts()
	if err != nil {
		panic(err)
	}
}

func loadLayouts() (map[Key]*Layout, error) {
	var err error
	once.Do(func() {
		entries, err := generatedLayouts.ReadDir(layouts)
		if err != nil {
			return
		}
		var i int
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			var data []byte
			data, err = generatedLayouts.ReadFile(filepath.Join(layouts, entry.Name()))
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
			structLayouts[key] = &lyt
			i++
		}
	})
	return structLayouts, err
}

// GetLayout returns the matching layout for the given version.
func GetLayout(v *semver.Version) (Key, *Layout, error) {
	layouts, err := loadLayouts()
	if err != nil {
		return Key{}, nil, err
	}
	for k, l := range layouts {
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
	layouts, err := loadLayouts()
	if err != nil {
		return nil, err
	}
	return layouts, nil
}

// loadInitialState loads the initial state for the supported versions.
func loadInitialState() (map[Key]*InitialState, error) {
	entries, err := generatedState.ReadDir(initialState)
	if err != nil {
		return nil, err
	}
	initialStates := map[Key]*InitialState{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := generatedState.ReadFile(filepath.Join(initialState, entry.Name()))
		if err != nil {
			return nil, err
		}
		ext := filepath.Ext(entry.Name())
		// Filter out non-yaml files.
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		var initState InitialState
		if err := yaml.Unmarshal(data, &initState); err != nil {
			return nil, err
		}
		rawConstraint := strings.TrimSuffix(entry.Name(), ext)
		constr, err := semver.NewConstraint(rawConstraint)
		if err != nil {
			return nil, err
		}
		key := Key{constraint: constr.String()}
		initialStates[key] = &initState
	}
	return initialStates, nil
}

// GetInitialState returns the initial state for the given version.
func GetInitialState(v *semver.Version) (Key, *InitialState, error) {
	state, err := loadInitialState()
	if err != nil {
		return Key{}, nil, err
	}
	for k, l := range state {
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

// GetInitialStates returns all the initial states for the supported versions.
func GetInitialStates() (map[Key]*InitialState, error) {
	return loadInitialState()
}
