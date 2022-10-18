# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

#!/usr/bin/env sh

mvn -B -DnewVersion=$(cat ./VERSION) -DgenerateBackupPoms=false versions:set
mvn package
