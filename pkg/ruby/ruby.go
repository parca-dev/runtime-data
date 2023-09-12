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
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed versions/*.yaml
var versions embed.FS

// GetVersions returns all the versions of Ruby that are supported.
func GetVersions() ([]RubyVersionOffsets, error) {
	entries, err := versions.ReadDir("versions")
	if err != nil {
		return nil, err
	}
	var versions []RubyVersionOffsets
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile("versions/" + entry.Name())
		if err != nil {
			return nil, err
		}
		var version RubyVersionOffsets
		err = yaml.Unmarshal(data, &version)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// GetVersionMap returns a map of Ruby version offsets.
func GetVersionMap() (map[string]RubyVersionOffsets, error) {
	versions, err := GetVersions()
	if err != nil {
		return nil, err
	}
	versionMap := make(map[string]RubyVersionOffsets)
	for _, pvo := range versions {
		version := fmt.Sprintf("%d.%d.%d", pvo.MajorVersion, pvo.MinorVersion, pvo.PatchVersion)
		versionMap[version] = pvo
	}
	return versionMap, nil
}