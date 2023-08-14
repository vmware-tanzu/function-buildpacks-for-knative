#!/bin/bash
# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

set -euo pipefail

srp provenance compile --saveto "${SRP_WORKING_DIR}"/prov3_fragment.json

srp provenance submit --verbose --path "${SRP_WORKING_DIR}"/prov3_fragment.json
