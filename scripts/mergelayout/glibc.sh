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

ARCH=${ARCH:-""}
target_archs=(
    amd64
    arm64
)
if [ -n "${ARCH}" ]; then
    target_archs=("${ARCH}")
fi

TARGET_DIR=${TARGET_DIR:-pkg/libc/glibc/layout}
rm -rf "${TARGET_DIR}" || true
mkdir -p "${TARGET_DIR}"
for arch in "${target_archs[@]}"; do
    mkdir -p "${TARGET_DIR}/${arch}"
    ./mergelayout -o "${TARGET_DIR}/${arch}" tmp/glibc/${arch}/layout/'glibc_*.yaml'
done
