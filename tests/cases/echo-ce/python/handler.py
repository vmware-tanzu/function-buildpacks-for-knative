# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

import logging
from typing import Any

def handler(attributes: Any, data: Any):
    payload = data
    # if 'ce-my-attr' not in attributes or attributes['ce-my-attr'] != "test":
    #     return attributes, "Incorrect 'my-header' value in request"
    return attributes, payload
