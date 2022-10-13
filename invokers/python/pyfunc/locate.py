# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

import importlib.util
import inspect
import pathlib
import sys
import typing
import types

import flask
import cloudevents.sdk.event.base as ce_sdk

from .config import Config

def find_func(search_path: str, module_name: str, function_name: str) -> typing.Callable:
    cfg = Config(search_path, module_name, function_name)
    
    workspace = pathlib.Path(cfg.search_path).resolve()

    for f in workspace.glob("*.py"):
        if f.stem == cfg.module_name:
            file = f
            break
    else:
        raise Exception(f"Module '{cfg.module_name}' not found in '{cfg.search_path}'")
    
    print(f"Importing from {file}")

    sys.path.insert(0, str(workspace))
    spec = importlib.util.spec_from_file_location(f.stem, f)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    func = _func_from_module(module, cfg.function_name)
    sys.path.pop(0)

    return func

def _func_from_module(module: types.ModuleType, handler_name: str) -> typing.Callable:
    funcs = []
    for (name, x) in inspect.getmembers(module, inspect.isfunction):
        if name != handler_name:
            continue
        sig = inspect.signature(x)
        print(f">>{name}: {sig}")
        for arg in sig.parameters.values():
            convert = ArgumentConversion(arg)
            if not convert.valid:
                break
        else:
            print(f">>> Matched sig {sig}")
            funcs.append(x)
    
    if len(funcs) == 0:
        raise Exception(f"Function '{handler_name}' not found in module '{module.__name__}'")

    return funcs[0]

class ArgumentConversion:
    def __init__(self, p: inspect.Parameter):
        self.name = p.name
        self._convert = None
        self.need_event = False
        self.unknownArg = None
        TYPE_TO_TRANSLATION = {
            ce_sdk.BaseEvent: (lambda x: x, True),
            flask.Request: (lambda x: x, False),
        }
        NAME_TO_TRANSLATION = {
            "event": (lambda x: x, True),
            "data": (lambda x: x.data, True),
            "payload": (lambda x: x.data, True),
            "attributes": (lambda x: {k: x[k] for k in x}, True),
            "req": (lambda x: x, False),
            "request": (lambda x: x, False),
            "body": (lambda x: x.get_data(), False),
            "headers": (lambda x: x.headers, False),
        }
        if p.annotation in TYPE_TO_TRANSLATION:
            self._convert, self.need_event = TYPE_TO_TRANSLATION[p.annotation]
        if p.name in NAME_TO_TRANSLATION:
            self._convert, self.need_event = NAME_TO_TRANSLATION[p.name]

        if self._convert is None and p.default == inspect.Parameter.empty:
            if p.kind not in (
                inspect.Parameter.VAR_POSITIONAL,
                inspect.Parameter.VAR_KEYWORD,
            ):
                self.unknownArg = p

    def convert(self, req: flask.Request, ce: ce_sdk.BaseEvent = None) -> typing.Any:
        if not self.valid:
            raise ValueError(f"Unable to convert {self.p} to a function argument.")
        if self.need_event:
            return self._convert(ce)
        return self._convert(req)

    @property
    def valid(self) -> bool:
        return self.unknownArg is None
