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
PACKAGE_DIR=${PACKAGE_DIR:-${TEMP_DIR}/apk}
BIN_DIR=${BIN_DIR:-${TEMP_DIR}/bin}
PACKAGE_NAME=${PACKAGE_NAME:-musl}
DEBUGINFO_DIR=${DEBUGINFO_DIR:-${TEMP_DIR}/debuginfo}

# https://dl-cdn.alpinelinux.org/alpine/
alpine_versions=(
    v3.0
    v3.1
    v3.2
    v3.3
    v3.4
    v3.5
    v3.6
    v3.7
    v3.8
    v3.9
    v3.10
    v3.11
    v3.12
    v3.13
    v3.14
    v3.15
    v3.16
    v3.17
    v3.18
    v3.19
)

convertArch() {
    case $1 in
    amd64)
        echo "x86_64"
        ;;
    arm64)
        echo "aarch64"
        ;;
    esac
}

echo "Downloading alpine runtimes"
for version in "${alpine_versions[@]}"; do
    ./apkdownload -u "https://dl-cdn.alpinelinux.org/alpine/${version}/main" -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${PACKAGE_NAME}"
done

echo "Extracting debuginfo from $BIN_DIR/$PACKAGE_NAME"
for arch in $BIN_DIR/$PACKAGE_NAME/*; do
    for version in $arch/*; do
        for variant in $version/*; do
            if [ -d "$variant" ]; then
                a=$(basename "$arch")
                v=$(basename "$version")
                linuxArch=$(convertArch "$(basename "$arch")")
                if [ $(basename "$variant") == "dbg" ]; then
                    dbginfo="$variant"/usr/lib/debug/lib/ld-musl-"$linuxArch".so.1.debug
                    if [ -f "$dbginfo" ]; then
                        echo "copying $dbginfo to $DEBUGINFO_DIR/$PACKAGE_NAME/$a/$v/"
                        dbgTarget="$DEBUGINFO_DIR"/$PACKAGE_NAME/$a/$v/
                        mkdir -p "$dbgTarget"
                        cp "$dbginfo" "$dbgTarget"
                        continue
                    fi
                fi
            fi
        done
    done
done
