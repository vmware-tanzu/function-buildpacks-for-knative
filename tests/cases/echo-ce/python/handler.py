import logging
from typing import Any

def handler(data: Any, attributes: dict):
    # Your function implementation goes here
    print("I'M HERE IM THE HANDLER")
    return attributes, "Hello World!"

# def handler(attributes: Any, data: Any):
#     payload = data
#     if 'my-header' not in attributes or attributes['my-header'] != "test":
#         return attributes, "Incorrect 'my-header' value in request"
#     return attributes, payload
