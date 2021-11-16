from pyfunc import find_func

import pytest
import os
import typing

function_asset_base_dir = os.path.join("tests", "functions")

def test_find_func():
    handler_handler = Helpers.handler_path("handler-handler")
    func = find_func(handler_handler)
    assert isinstance(func, typing.Callable)

def test_find_func_specify_handler():
    os.environ["PYTHON_HANDLER"] = "handler.bar"
    handler_bar = Helpers.handler_path("handler-bar")
    func = find_func(handler_bar)
    assert isinstance(func, typing.Callable)

def test_find_func_invalid_handler():
    os.environ["PYTHON_HANDLER"] = "handler.unknown"
    handler_bar = Helpers.handler_path("handler-handler")
    with pytest.raises(Exception) as ex:
        find_func(handler_bar)
    assert "not found in module" in str(ex)

class Helpers:
    @staticmethod
    def handler_path(case):
        return os.path.join(function_asset_base_dir, "http", case)
