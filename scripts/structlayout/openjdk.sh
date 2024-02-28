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
SOURCE_DIR=${SOURCE_DIR:-${TEMP_DIR}/debuginfo}
TARGET_DIR=${TARGET_DIR:-${TEMP_DIR}/openjdk}

mkdir -p "${TARGET_DIR}"
for repo in $(ls -d ${SOURCE_DIR}/openjdk-*); do
    for arch in $(ls -d ${repo}/*); do
        for version in $(ls -d ${arch}/*); do
            for dbgfile in $(ls -d ${version}/*); do
                v=$(basename "${version}")
                a=$(basename "${arch}")
                echo "Running structlayout against openjdk ${v} for ${a}..."
                ./structlayout -r java -v "${v}" -o "${TARGET_DIR}/${a}" "${dbgfile}"
            done
        done
    done
done
