import logging
from typing import Any

def handler(data: Any, attributes: Any):
    payload = data
    if 'my-header' not in attributes or attributes['my-header'] != "test":
        return attributes, "Incorrect 'my-header' value in request"
    return attributes, payload
