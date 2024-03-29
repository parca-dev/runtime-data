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

# This script helps to download python runtimes using container images for integration tests.

CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}
TARGET_DIR=${TARGET_DIR:-tests/integration/binaries/python}
ARCH=${ARCH:-""}

# https://devguide.python.org/versions/
python_versions=(
    # Unsupported Versions (end of life)
    2.7.15
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
    3.12.0
    3.13.0a4
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

# Install libpython for each python version under python_runtimes directory.
for python_version in "${python_versions[@]}"; do
    for arch in "${target_archs[@]}"; do
        target="${PWD}"/"${TARGET_DIR}"/"${arch}"/${python_version}
        echo "Checking if python ${python_version} runtime for ${arch} is already downloaded..."
        if ls "${target}"/libpython"${python_version%.*}"*.so.1.0 1>/dev/null 2>&1; then
            echo "Python ${python_version} runtime for ${arch} is already downloaded."
            continue
        fi
        echo "Downloading python ${python_version} runtime..."
        "${CONTAINER_RUNTIME}" run --rm \
            --platform "linux/${arch}" \
            -v "${target}":/tmp -w /tmp \
            python:${python_version} \
            bash -c 'cp /usr/local/lib/libpython"${python_version%.*}"*.so.1.0 \
            /tmp'
        echo "Changing the owner of the file to the current user..."
        sudo chown -R $(whoami) "${TARGET_DIR}"
        echo "Done."
        for i in "${target}"/libpython*; do file "$i"; done
    done
done

echo "All python runtimes downloaded successfully."
dir="${PWD}"/"${TARGET_DIR}"/*/*/libpython*
for i in $dir; do file "$i"; done
