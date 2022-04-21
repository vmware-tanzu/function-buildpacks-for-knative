# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from functools import reduce
from flask_healthz import healthz
from flask_healthz import HealthError
import inspect
import os
import typing
import traceback

import flask
import cloudevents.http
from cloudevents.http.util import default_marshaller
from cloudevents.sdk import types

from .locate import ArgumentConversion, find_func

# There is a bug in the cloudevent sdk where if the data contents is a string it will
# run it through the json marshaller which ends up wrapping it with double quotes.
# The offending code lives here: https://github.com/cloudevents/sdk-python/blob/master/cloudevents/sdk/event/base.py#L285-L312
def smart_marshaller(content: typing.Any):
    if isinstance(content, str):
        return content
    else:
        return default_marshaller(content)

def WrapFunction(func: typing.Callable) -> typing.Callable:
    sig = inspect.signature(func)
    args = [ArgumentConversion(p) for p in sig.parameters.values()]
    need_cloudevent = reduce(lambda a, b: a or b, (p.need_event for p in args))

    def handler() -> flask.Response:
        req = flask.request
        event = None
        if need_cloudevent:
            event = cloudevents.http.from_http(req.headers, req.get_data())

        params = {arg.name: arg.convert(req, event) for arg in args}

        result = func(**params)

        if isinstance(result, flask.Response):
            return result

        if need_cloudevent and result:
            if not isinstance(result, cloudevents.http.CloudEvent):
                try:
                    if isinstance(result, tuple):
                        result = cloudevents.http.CloudEvent(*result)
                    else:
                        result = cloudevents.http.CloudEvent(result)
                except Exception:
                    # Opportunistic conversion; failures will be caught in the next block
                    pass
            try:
                headers, body = cloudevents.http.to_binary(result, data_marshaller=smart_marshaller)
                return flask.Response(body, 200, headers)
            except Exception as e:
                print(f"Sending raw result: {e}")
                traceback.print_exc()
                return flask.Response(
                    "Accepted with no event response", 200, mimetype="text/plain"
                )

        return flask.Response(result, 200)

    print(f"$$ Converting {inspect.signature(func)} to {inspect.signature(handler)}")
    return handler


def main(dir: str = "."):
    func = find_func(dir)
    http_func = WrapFunction(func)
    # TODO: add option for GET / handle multiple functions
    app = flask.Flask(func.__name__)
    app.register_blueprint(healthz, url_prefix="/")
    app.config['']
    app.add_url_rule("/", view_func=http_func, methods=["POST","GET"])
    app.run(host="0.0.0.0", port=os.environ.get("PORT", 8080))
