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

TMP_DIR=${TMP_DIR:-tmp}
OUTPUT_DIR=${OUTPUT_DIR:-${TMP_DIR}/python}
INPUT_DIR=${INPUT_DIR:-tests/integration/binaries/python}

mkdir -p tmp/python/
for arch in "${INPUT_DIR}"/*; do
    for version in "${arch}"/*; do
        arch=$(basename "${arch}")
        version=$(basename "${version}")
        echo "Running structlayout against python ${version} runtime for ${arch}..."
        ./structlayout -r python -v "${version}" -o ${OUTPUT_DIR}/${arch} ${INPUT_DIR}/${arch}/"${version}"/libpython"${version%.*}"*.so.1.0
    done
done
