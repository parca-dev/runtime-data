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
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed versions/*.yaml
var versions embed.FS

// GetVersions returns all the versions of Python that are supported.
func GetVersions() ([]PythonVersionOffsets, error) {
	entries, err := versions.ReadDir("versions")
	if err != nil {
		return nil, err
	}
	var versions []PythonVersionOffsets
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile("versions/" + entry.Name())
		if err != nil {
			return nil, err
		}
		var version PythonVersionOffsets
		err = yaml.Unmarshal(data, &version)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// GetVersionMap returns a map of Python version offsets.
func GetVersionMap() (map[string]PythonVersionOffsets, error) {
	versions, err := GetVersions()
	if err != nil {
		return nil, err
	}
	versionMap := make(map[string]PythonVersionOffsets)
	for _, pvo := range versions {
		version := fmt.Sprintf("%d.%d.%d", pvo.MajorVersion, pvo.MinorVersion, pvo.PatchVersion)
		versionMap[version] = pvo
	}
	return versionMap, nil
}
