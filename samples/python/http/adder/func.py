# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from typing import Any
import flask

# This is a toy sample and should not be used as an example of how to authenticate a user.
# Credentials should never be stored like this and errors returned should not be so specific.

creds = {"admin": "supersecure", "asu": "mypassword", "someone": "123qwe"}

def main(req: Any):
    username = req.form.get('username')
    if username not in creds:
        return flask.Response("Unauthorized: User does not exist", 401)
    
    password = req.form.get("password")
    if password != creds[username]:
        return flask.Response("Unauthorized: Wrong password", 401) # Return your own flask response to set a return code!

    first = req.form.get("first", default=0)
    second = req.form.get("second", default=0)

    sum = int(first) + int(second)
    return f"Hello, {username}! The answer to {first} + {second} is {sum}\n" # Without returning a flask response, it'll automatically be HTTPOK
    
