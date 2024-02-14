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

# This script helps to download ruby runtimes using container images for integration tests.

CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}
TARGET_DIR=${TARGET_DIR:-tests/integration/binaries/ruby}
ARCH=${ARCH}

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

# Install libpython for each ruby version under python_runtimes directory.
for ruby_version in "${ruby_versions[@]}"; do
    for arch in "${target_archs[@]}"; do
        target="${PWD}"/"${TARGET_DIR}"/"${arch}"
        echo "Checking if ruby ${ruby_version} runtime for ${arch} is already downloaded..."
        if ls "${target}"/libruby.so.${ruby_version} 1>/dev/null 2>&1; then
            echo "Ruby ${ruby_version} runtime for ${arch} is already downloaded."
            continue
        fi
        echo "Downloading ruby ${ruby_version} runtime..."
        "${CONTAINER_RUNTIME}" run --rm \
            --platform "linux/${arch}" \
            -v "${target}":/tmp -w /tmp \
            ruby:${ruby_version}-slim \
            cp /usr/local/lib/libruby.so.${ruby_version} \
            /tmp
        echo "Changing the owner of the file to the current user..."
        sudo chown -R $(whoami) "${TARGET_DIR}"
        echo "Done."

    done
done

echo "All ruby runtimes downloaded successfully."
dir="${PWD}"/"${TARGET_DIR}"/**/libruby*
for i in $dir; do file "$i"; done
