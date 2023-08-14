#!/bin/bash
# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

echo "SRP: Collect source provenance."

mkdir -p ${SRP_WORKING_DIR}
srp config auth --client-id "$CLIENT_ID" --client-secret "$CLIENT_SECRET"
srp provenance init

GITHUB_FQDN=$(echo "${GITHUB_SERVER_URL}" | sed -e "s/^https:\/\///")
srp provenance add-build github --action ${GITHUB_ACTION} --build-id ${GITHUB_RUN_ID}_${GITHUB_RUN_ATTEMPT} --instance ${GITHUB_FQDN} --namespace ${GITHUB_REPOSITORY} --ref ${GITHUB_REF}
srp provenance declare-source git --verbose --set-key=function-buildpack-source --path .
srp provenance action start --name=publish

srp provenance action import-cmd --cmd "make base_url=$url registry.location=other REGISTRY=$registry $target"
observer_agent -m start_observer -e "${SRP_WORKING_DIR}"/required-envs.sh -S
source "${SRP_WORKING_DIR}"/required-envs.sh set

make base_url=$url registry.location=other REGISTRY=$registry $target

source "${SRP_WORKING_DIR}"/required-envs.sh unset
rm "${SRP_WORKING_DIR}/required-envs.sh"
observer_agent -m stop_observer -f network_provenance.json

key="${registry}/${BUILDPACK}-buildpack:${version}"
echo "key set to ${key}"
action="publish"
image=$(docker inspect "${key}" |  jq -r '.[0].RepoDigests[0]')
srp provenance add-output docker \
    	--set-key="${key}" \
    	--action-key="${action}" \
    	--name="${image%%@*}" \
    	--location="${image%%@*}" \
    	--digest="${image##*@}"

srp provenance add-input syft --scan-target="${key}" --output-key="${key}"

srp provenance action import-observation --name="publish" --file=network_provenance.json
srp provenance action stop

echo "------------- Completed ${SRP_WORKING_DIR}/_provenance.json -------------"
cat "${SRP_WORKING_DIR}"/_provenance.json
