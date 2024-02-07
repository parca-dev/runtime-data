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

# This script helps to download python runtimes using container images.

CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}

# https://devguide.python.org/versions/
python_versions=(
    # Unsupported Versions (end of life)
    2.7.15
    3.3.7
    3.4.8
    3.5.5
    3.6.6
    3.7.0
    # Supported Versions
    3.8.0  # security
    3.9.5  # security
    3.10.0 # security
    3.11.0 # bugfix
)

# Check if CONTAINER_RUNTIME is installed.
if ! command -v "${CONTAINER_RUNTIME}" &>/dev/null; then
    echo "ERROR: ${CONTAINER_RUNTIME} is not installed."
    exit 1
fi

# Install libpython for each python version under python_runtimes directory.
for python_version in "${python_versions[@]}"; do
    echo "Downloading python ${python_version} runtime..."
    "${CONTAINER_RUNTIME}" run --rm -v "${PWD}"/tests/integration/tmp:/tmp -w /tmp python:${python_version} bash -c 'cp /usr/local/lib/libpython"${python_version%.*}"*.so.1.0 /tmp'
    echo "Done."
done

echo "All python runtimes downloaded successfully."

for i in tests/integration/tmp/libpython*; do file "$i"; done
