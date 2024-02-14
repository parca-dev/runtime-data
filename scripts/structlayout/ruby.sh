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

# This script helps to run structlayout for specified ruby layout for integration tests.

ruby_versions=(
    2.6.0
    2.6.3
    2.7.1
    2.7.4
    2.7.6
    3.0.0
    3.0.4
    3.1.2
    3.1.3
    3.2.0
    3.2.1
)

mkdir -p tmp/ruby
for ruby_version in "${ruby_versions[@]}"; do
    echo "Running structlayout for ruby ${ruby_version} runtime..."
    ./structlayout -r ruby -v "${ruby_version}" -o tmp/ruby tests/integration/binaries/ruby/libruby.so.${ruby_version}
done