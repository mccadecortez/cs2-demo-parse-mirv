#!/usr/bin/env bash

# XXX:: dirname twice to get project root
export PROJECT="submodules/mirv-script"
export SCRIPT_DIR=$(dirname $(dirname $(realpath -s $0)))
export PROJECT_DIR="${SCRIPT_DIR}/${PROJECT}"

if [[ ! "${SCRIPT_DIR+x}" ]]; then
    echo "SCRIPT_DIR is empty (bug)"

    exit 1
fi

if ! command -v npm &>/dev/null; then
    echo "missing 'npm' to build project ${PROJECT_DIR}"

    exit 1
fi

## Build
set -e -x

cd "${PROJECT_DIR}"
npm run build
npm run build-scripts

echo "${PROJECT_DIR}"
exit 0