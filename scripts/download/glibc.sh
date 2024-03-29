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

TEMP_DIR=${TEMP_DIR:-tmp}
PACKAGE_DIR=${PACKAGE_DIR:-${TEMP_DIR}/deb}
BIN_DIR=${BIN_DIR:-${TEMP_DIR}/bin}
DEBUGINFO_DIR=${DEBUGINFO_DIR:-${TEMP_DIR}/debuginfo}
PACKAGE_NAME=${PACKAGE_NAME:-libc6}

./debdownload -u 'http://ftp.debian.org/debian/pool/main/g/glibc/' -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${PACKAGE_NAME}"
./debdownload -u 'http://archive.ubuntu.com/ubuntu/pool/main/g/glibc/' -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${PACKAGE_NAME}"
./debdownload -u 'http://old-releases.ubuntu.com/ubuntu/pool/main/g/glibc/' -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${PACKAGE_NAME}"

# 2.35 for arm64 is not available the archives above, so we need to download it separately.
./debdownload -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${PACKAGE_NAME}" -s http://launchpadlibrarian.net/707340001/libc6_2.35-0ubuntu3.6_arm64.deb
./debdownload -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${PACKAGE_NAME}" -s http://launchpadlibrarian.net/707339998/libc6-dbg_2.35-0ubuntu3.6_arm64.deb

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

echo "Copying debuginfo from $BIN_DIR/$PACKAGE_NAME"
for arch in $BIN_DIR/$PACKAGE_NAME/*; do
    for version in $arch/*; do
        for variant in $version/*; do
            if [ -d "$variant" ]; then
                a=$(basename "$arch")
                v=$(basename "$version")
                linuxArch=$(convertArch "$(basename "$arch")")
                if [ $(basename "$variant") == "dbg" ]; then
                    if [ -d "$variant"/usr/lib/debug/lib ]; then
                        # ubuntu
                        dbginfo="$variant"/usr/lib/debug/lib/"$linuxArch"-linux-gnu/libc-"${v%.*}".so
                        if [ ! -f "$dbginfo" ]; then
                            dbginfo="$variant"/usr/lib/debug/lib/"$linuxArch"-linux-gnu/libc.so.6
                        fi
                        if [ -f "$dbginfo" ]; then
                            echo "Copying $dbginfo to $DEBUGINFO_DIR/$PACKAGE_NAME/$a/$v/"
                            dbgTarget="$DEBUGINFO_DIR"/$PACKAGE_NAME/$a/$v/
                            mkdir -p "$dbgTarget"
                            cp "$dbginfo" "$dbgTarget"
                            continue
                        fi
                    fi
                    # debian
                    continue
                fi
                if [ $(basename "$variant") == "main" ]; then
                    target="$variant"/lib/"$linuxArch"-linux-gnu/libc-"${v%.*}".so
                    if [ ! -f "$target" ]; then
                        target="$variant"/lib/"$linuxArch"-linux-gnu/libc.so.6
                    fi
                    if [ ! -f "$target" ]; then
                        target="$variant"/usr/lib/"$linuxArch"-linux-gnu/libc.so.6
                    fi
                    dbginfo=$(./debuginfofind -d "$version"/dbg $target || true)
                    if [ -n "$dbginfo" ]; then
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
