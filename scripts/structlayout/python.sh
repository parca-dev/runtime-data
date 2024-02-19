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

# This script helps to run structlayout for specified python versions for integration tests.

# https://devguide.python.org/versions/
python_versions=(
    # Unsupported Versions (end of life)
    2.7.15
    2.7.18
    3.3.7
    3.4.8
    3.5.5
    3.6.6
    3.7.0
    3.7.2
    3.7.4
    # Supported Versions
    3.8.0  # security
    3.8.1  # security
    3.9.5  # security
    3.9.6  # security
    3.10.0 # security
    3.11.0 # bugfix
)

ARCH=${ARCH:-""}
target_archs=(
    amd64
    arm64
)
if [ -n "${ARCH}" ]; then
    target_archs=("${ARCH}")
fi

mkdir -p tmp/python/
for python_version in "${python_versions[@]}"; do
    for arch in "${target_archs[@]}"; do
        echo "Running structlayout against python ${python_version} runtime for ${arch}..."
        ./structlayout -r python -v "${python_version}" -o tmp/python/${arch} tests/integration/binaries/python/${arch}/"${python_version}"/libpython"${python_version%.*}"*.so.1.0
    done
done
