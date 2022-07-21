#!/bin/bash
# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

trap 'catch' ERR
catch() {
  echo "An error has occurred removing SRP data"
  rm -rf ./srp_data
}

if [ -n "${DOMAIN}" ]; then
  DOMAIN="domain='${DOMAIN}',"
fi

GITHUB_FQDN=$(echo "${GITHUB_SERVER_URL}" | sed -e "s/^https:\/\///")
SRP_UID="uid.obj.build.github(${DOMAIN}instance='${GITHUB_FQDN}',namespace='${GITHUB_REPOSITORY}',ref='${GITHUB_REF}',action='${GITHUB_ACTION}',build_id='${GITHUB_RUN_ID}_${GITHUB_RUN_ATTEMPT}')"
echo "SRP component UID generated: $SRP_UID"

mkdir -p srp_data/

echo "$SRP_UID" > srp_data/srp_uid
echo "SRP component UID stored in: srp_data/srp_uid"

echo "${GITHUB_RUN_ID}_${GITHUB_RUN_ATTEMPT}" > srp_data/build_number
echo "Build number stored in:      srp_data/build_number"

srp config auth --client-id "$CLIENT_ID" --client-secret "$CLIENT_SECRET"
cp "$HOME/.srp/config.yml" srp_data/config.yml
echo "SRP CLI config stored in:    srp_data/config.yml"
echo "SRP CLI version:             $(srp --version)"

echo "::set-output name=srp-data::$PWD/srp_data"
