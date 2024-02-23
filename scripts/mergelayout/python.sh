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

TEMP_DIR=${TEMP_DIR:-tmp}
OUTPUT_DIR=${OUTPUT_DIR:-pkg/python}
INPUT_DIR=${INPUT_DIR:-${TEMP_DIR}/python}

rm -rf "${OUTPUT_DIR}"/layout
rm -rf "${OUTPUT_DIR}"/initialstate
for arch in "${INPUT_DIR}"/*; do
    arch=$(basename "${arch}")
    mkdir -p "${OUTPUT_DIR}"/layout/"${arch}"
    ./mergelayout -o "${OUTPUT_DIR}"/layout/"${arch}" "${INPUT_DIR}"/"${arch}"/layout/'python_*.yaml'

    mkdir -p "${OUTPUT_DIR}"/initialstate"/${arch}"
    ./mergelayout -o "${OUTPUT_DIR}"/initialstate/"${arch}" "${INPUT_DIR}"/"${arch}"/initialstate/'python_*.yaml'
done
