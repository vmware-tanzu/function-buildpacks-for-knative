#!/bin/bash
# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

echo "Installing SRP CLI"

mkdir -p "$HOME/srp-tools"
echo "$HOME/srp-tools" >> "$GITHUB_PATH"
echo "$HOME/srp-tools/observer/bin" >> "$GITHUB_PATH"

SRP_CLI_VERSION='0.9.9-20230724044630-61ef470-169'
curl -L "${SRP_TOOLS_URL}/${SRP_CLI_VERSION}/srp-tools-linux-amd64-${SRP_CLI_VERSION}.tar.gz" | tar -xz -C "$HOME/srp-tools"

"$HOME/srp-tools/srp" --version
"$HOME/srp-tools/observer/bin/observer_agent" --version
