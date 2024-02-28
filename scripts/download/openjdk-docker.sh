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
TMP_DIR=${TMP_DIR:-tmp}
TARGET_DIR=${TARGET_DIR:-${TMP_DIR}/openjdk/binaries}
ARCH=${ARCH:-""}

openjdk_images=(
    eclipse-temurin:17
    amazoncorretto:17
    sapmachine:17

    eclipse-temurin:18
    amazoncorretto:18

    eclipse-temurin:19
    amazoncorretto:19
    sapmachine:19

    eclipse-temurin:20
    amazoncorretto:20
    sapmachine:20

    eclipse-temurin:21
    amazoncorretto:21
    sapmachine:21
)
target_archs=(
    amd64
    arm64
)

# Check if CONTAINER_RUNTIME is installed.
if ! command -v "${CONTAINER_RUNTIME}" &>/dev/null; then
    echo "ERROR: ${CONTAINER_RUNTIME} is not installed."
    exit 1
fi

# TODO(kakkoyun): We need to install debug info packages.
# Install libjvm.so for each openjdk image.
for openjdk_image in "${openjdk_images[@]}"; do
    for target_arch in "${target_archs[@]}"; do
        image_name=${openjdk_image/:/-}
        version=${openjdk_image#*:}
        target_dir="${PWD}/${TARGET_DIR}/${image_name}/${target_arch}"
        # Skip if libjvm.so already exists.
        if [ -f "${target_dir}/libjvm.so" ]; then
            continue
        fi
        mkdir -p "${target_dir}"
        echo "Downloading ${openjdk_image} for ${target_arch} to ${target_dir}..."
        sources=(
            /opt/java/openjdk/lib/server/libjvm.so
            /usr/lib/jvm/java/lib/server/libjvm.so
            /usr/lib/jvm/sapmachine"-${version}"/lib/server/libjvm.so
        )
        for source in "${sources[@]}"; do
            if ${CONTAINER_RUNTIME} run --rm \
                --platform "linux/${target_arch}" \
                -v "${target_dir}":/tmp \
                -w /tmp \
                "${openjdk_image}" \
                bash -c "cp ${source} /tmp/"; then
                break
            fi
        done
    done
done

echo "All openjdk runtimes downloaded successfully."
dir="${PWD}/${TARGET_DIR}"/*/*/libjvm.so
for i in $dir; do file "$i"; done
