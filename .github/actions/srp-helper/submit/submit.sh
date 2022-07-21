#!/bin/bash
# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

trap 'catch' ERR
catch() {
  echo "An error has occurred removing SRP data"
  rm -rf ./srp_data
}

echo "Runway SRP: Submit provenance."

# TODO: remove sed call
#  At the moment, the SRP CLI does not properly URL encode the UID. This causes
#  the API call to fail. The SRP CLI released in a few weeks should fix this.
#  When that happens, it'll be safe to remove the sed translation.
SRP_UID="$(sed 's|/|%2F|g' < ./srp_data/srp_uid)"
FULL_SRP_UID="uid.mtd.provenance_2_5.fragment(obj_uid=$SRP_UID,revision='')"
echo "Full SRP UID that will be used for upload: $FULL_SRP_UID"

if [ -z "$SOURCE_PROVENANCE_PATH" ] && [ -f "$SOURCE_PROVENANCE_PATH" ]; then
    cp "$SOURCE_PROVENANCE_PATH" "./srp_data/source_provenance.json"
fi

if [ -z "$NETWORK_PROVENANCE_PATH" ] && [ -f "$NETWORK_PROVENANCE_PATH" ]; then
    cp "$NETWORK_PROVENANCE_PATH" "./srp_data/network_provenance.json"
fi

if [ -f "./srp_data/network_provenance.json" ]; then
    echo "Found network provenance: ./srp_data/network_provenance.json"
    echo "Merging network provenance and source provenance:"
    srp provenance merge \
        --source ./srp_data/source_provenance.json \
        --network ./srp_data/network_provenance.json \
        --saveto ./srp_data/finalized_source_provenance.json \
        --config ./srp_data/config.yml
else
    echo "No network provenance found"
    cp ./srp_data/source_provenance.json ./srp_data/finalized_source_provenance.json
fi

# TODO: move --url to the init phase, once the SRP CLI supports it
#   There is an open issue to allow it to be set through config.yml, which would 
#   remove the need to set it on every invocation.
echo "Finalized source provenance location: ./srp_data/finalized_source_provenance.json"
echo "Submitting source provenance via SRP CLI:"
if [ -z "$SRP_URL" ]; then
    srp metadata submit \
        --path ./srp_data/finalized_source_provenance.json \
        --uid "$FULL_SRP_UID" \
        --config ./srp_data/config.yml
else
    srp metadata submit \
        --path ./srp_data/finalized_source_provenance.json \
        --uid "$FULL_SRP_UID" \
        --url "$SRP_URL" \
        --config ./srp_data/config.yml
fi

echo "Downloading source provenance via SRP CLI:"
if [ -z "$SRP_URL" ]; then
    srp metadata get \
        --uid "$FULL_SRP_UID" \
        --config ./srp_data/config.yml \
        --pretty
else
    srp metadata get \
        --uid "$FULL_SRP_UID" \
        --config ./srp_data/config.yml \
        --url "$SRP_URL" \
        --pretty
fi

echo "Done. Removing SRP data."
rm -rf ./srp_data
