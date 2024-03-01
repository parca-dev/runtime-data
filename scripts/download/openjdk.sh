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
PACKAGE_DIR=${PACKAGE_DIR:-${TEMP_DIR}/deb}
BIN_DIR=${BIN_DIR:-${TEMP_DIR}/bin}
DEBUGINFO_DIR=${DEBUGINFO_DIR:-${TEMP_DIR}/debuginfo}

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

# Download openjdk with debug info.
repositories=(
    "openjdk-17"
    "openjdk-19"
    "openjdk-20"
    "openjdk-21"
    "openjdk-22"
    "openjdk-23"
)

for repo in "${repositories[@]}"; do
    ./debdownload -u "http://ftp.debian.org/debian/pool/main/o/${repo}/" -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${repo}" -v "jre-headless" || true
    ./debdownload -u "http://archive.ubuntu.com/ubuntu/pool/main/o/${repo}/" -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${repo}" -v "jre-headless" || true
    ./debdownload -u "http://old-releases.ubuntu.com/ubuntu/pool/main/o/${repo}/" -t "${PACKAGE_DIR}" -o "${BIN_DIR}" -p "${repo}" -v "jre-headless" || true
done

for repo in "${repositories[@]}"; do
    echo "Extracting debuginfo from $BIN_DIR/$repo"
    for arch in $BIN_DIR/$repo/*; do
        for version in $arch/*; do
            for variant in $version/*; do
                if [ -d "$variant" ]; then
                    a=$(basename "$arch")
                    v=$(basename "$version")
                    if [ $(basename "$variant") == "jre-headless" ]; then
                        jvm_version=${repo#*-}
                        # usr/lib/jvm/java-17-openjdk-amd64/lib/server/libjvm.so
                        target="$variant"/usr/lib/jvm/java-"$jvm_version"-openjdk-"$a"/lib/server/libjvm.so
                        dbginfo=$(./debuginfofind -d "$version"/dbg $target || true)
                        if [ -f "$dbginfo" ]; then
                            echo "Copying $dbginfo to $DEBUGINFO_DIR/$repo/$a/$v/"
                            dbgTarget="$DEBUGINFO_DIR"/$repo/$a/$v/
                            mkdir -p "$dbgTarget"
                            cp "$dbginfo" "$dbgTarget"
                        fi
                    fi
                fi
            done
        done
    done
done

echo "All debuginfo extracted to $DEBUGINFO_DIR"
# /home/kakkoyun/Workspace/Projects/parca/runtime-data/tmp/debuginfo/openjdk-17/amd64/17.0.10/5835b4c09e3197e1f64d9fd0baf17e665004c8.debug
glob="$PWD"/"$DEBUGINFO_DIR"/openjdk-*/*/*/*.debug
for f in $glob; do file "$f"; done
