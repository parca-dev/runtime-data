#!/usr/bin/env bash

# Copyright 2022-2024 The Parca Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -euo pipefail

# This script helps to merge structlayout outputs in specified directory for integration tests.

target_archs=(
    amd64
    arm64
)

rm -rf pkg/ruby/layout
for arch in "${target_archs[@]}"; do
    mkdir -p pkg/ruby/layout/"${arch}"
    ./mergelayout -o pkg/ruby/layout/"${arch}" tmp/ruby/"${arch}"/layout/'ruby_*.yaml'
done
