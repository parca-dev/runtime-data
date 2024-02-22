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

TEMP_DIR=${TEMP_DIR:-tmp}
SOURCE_DIR=${SOURCE_DIR:-${TEMP_DIR}/debuginfo/libc6}
TARGET_DIR=${TARGET_DIR:-${TEMP_DIR}/glibc}

mkdir -p "${TARGET_DIR}"
for arch in "${SOURCE_DIR}"/*; do
    for version in "${arch}"/*; do
        for dbgfile in "${version}"/*; do
            v=$(basename "${version}")
            a=$(basename "${arch}")
            echo "Running structlayout against glibc ${v} for ${a}..."
            ./structlayout -r glibc -v "${v}" -o "${TARGET_DIR}/${a}" "${dbgfile}"
        done
    done
done
