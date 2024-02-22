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

CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}
ARCH=${ARCH:-""}

target_archs=(
    amd64
    arm64
)
if [ -n "${ARCH}" ]; then
    target_archs=("${ARCH}")
fi

# Check if CONTAINER_RUNTIME is installed.
if ! command -v "${CONTAINER_RUNTIME}" &>/dev/null; then
    echo "ERROR: ${CONTAINER_RUNTIME} is not installed."
    exit 1
fi

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

# Check libc versions for each python version.
for python_version in "${python_versions[@]}"; do
    for arch in "${target_archs[@]}"; do
        "${CONTAINER_RUNTIME}" run --rm \
            --platform "linux/${arch}" \
            python:${python_version} \
            bash -c 'ldd -r -v /usr/local/lib/libpython"${python_version%.*}"*.so.1.0'
    done
done

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

# Check glibc versions for each given ruby version.
for ruby_version in "${ruby_versions[@]}"; do
    for arch in "${target_archs[@]}"; do
        "${CONTAINER_RUNTIME}" run --rm \
            --platform "linux/${arch}" \
            ruby:${ruby_version}-slim \
            ldd -r -v /usr/local/lib/libruby.so."${ruby_version}"
    done
done
