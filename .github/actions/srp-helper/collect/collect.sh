#!/bin/bash
# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

trap 'catch' ERR
catch() {
  echo "An error has occurred removing SRP data"
  rm -rf ./srp_data
}

echo "Runway SRP: Collect source provenance."

if [ ! -d "./srp_data" ]; then
    echo 'SRP Data directory missing, you must run "init" first!'
    exit 1
fi

SRP_UID="$(xargs echo -n < ./srp_data/srp_uid)"
BUILD_NUM="$(xargs echo -n < ./srp_data/build_number)"

for REPO in $GIT_REPOS; do
    echo "Collecting source provenance for '$REPO' via SRP CLI"
    if [ -f "./srp_data/source_provenance.json" ]; then
        srp provenance source \
            --path "$REPO" \
            --comp-uid "$SRP_UID" \
            --all-ephemeral true \
            --scm-type "$SCM_TYPE" \
            --build-type "$BUILD_TYPE" \
            --build-number "$BUILD_NUM" \
            --saveto "./srp_data/source_provenance.json" \
            --append
    else
        srp provenance source \
            --path "$REPO" \
            --comp-uid "$SRP_UID" \
            --all-ephemeral true \
            --scm-type "$SCM_TYPE" \
            --build-type "$BUILD_TYPE" \
            --build-number "$BUILD_NUM" \
            --saveto "./srp_data/source_provenance.json"
    fi
done

echo "Provenance stored in: ./srp_data/source_provenance.json"
echo "Collected provenance:"
cat ./srp_data/source_provenance.json
