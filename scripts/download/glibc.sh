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

# This script helps to download glibc using debian packages.

PACKAGE_DIR=${PACKAGE_DIR:-tmp/deb}
BIN_DIR=${BIN_DIR:-tmp/bin}
DEBUGINFO_DIR=${DEBUGINFO_DIR:-tmp/debuginfo}
PACKAGE_NAME=${PACKAGE_NAME:-libc6}

./debdownload -u 'http://ftp.debian.org/debian/pool/main/g/glibc/' -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${PACKAGE_NAME}"

# TODO: Add ubuntu and ubuntu old releases.
# ./debdownload -u 'http://archive.ubuntu.com/ubuntu/pool/main/g/glibc/' -t tmp/deb -o tmp/bin -p 'libc6'
# ./debdownload -u 'http://old-releases.ubuntu.com/ubuntu/pool/main/g/glibc/' -t tmp/deb -o tmp/bin -p 'libc6'

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

for arch in $BIN_DIR/$PACKAGE_NAME/*; do
    for version in $arch/*; do
        for variant in $version/*; do
            if [ -d "$variant" ]; then
                if [ $(basename "$variant") == "main" ]; then
                    v=$(basename "$version")
                    linuxArch=$(convertArch "$(basename "$arch")")
                    target="$variant"/lib/"$linuxArch"-linux-gnu/libc-"${v%.*}".so
                    if [ ! -f "$target" ]; then
                        target="$variant"/lib/"$linuxArch"-linux-gnu/libc.so.6
                    fi
                    dbginfo=$(./debuginfofind -d "$version"/dbg $target)
                    if [ -n "$dbginfo" ]; then
                        a=$(basename "$arch")
                        v=$(basename "$version")
                        target="$DEBUGINFO_DIR"/$PACKAGE_NAME/$a/$v/
                        mkdir -p "$target"
                        cp "$dbginfo" "$target"
                    fi
                fi
            fi
        done
    done
done
