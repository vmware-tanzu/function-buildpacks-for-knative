# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

has_srp_creds=$(test "${CLIENT_ID}" &> /dev/null && test "${CLIENT_SECRET}" &> /dev/null && echo 'true')
echo "::set-output name=has-srp-creds::${has_srp_creds}"