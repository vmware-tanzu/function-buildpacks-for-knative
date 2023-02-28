#!/bin/bash
# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

echo "Installing SRP CLI"

mkdir -p "$HOME/bin"
echo "$HOME/bin" >> "$GITHUB_PATH"

curl \
  --show-error \
  --silent \
  --location \
  --output "$HOME/bin/srp" \
  "$SRP_CLIENT_URL"

chmod 755 "$HOME/bin/srp"

"$HOME/bin/srp" --version
