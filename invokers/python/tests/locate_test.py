from pyfunc import find_func
from pyfunc import constants

import pytest
import os
import typing
from unittest import mock

function_asset_base_dir = os.path.join("tests", "functions")

def test_find_func():
    path = Helpers.search_path('default')
    func = find_func(path)
    assert isinstance(func, typing.Callable)

def test_specify_module():
    with mock.patch.dict(os.environ, {constants.ENV_MODULE_NAME: 'handler'}):
        path = Helpers.search_path('module-handler')
        func = find_func(path)
        assert isinstance(func, typing.Callable)

def test_specify_function():
    with mock.patch.dict(os.environ, {constants.ENV_FUNCTION_NAME: 'foo'}):
        path = Helpers.search_path('func-foo')
        func = find_func(path)
        assert isinstance(func, typing.Callable)

def test_module_not_found():
    module = 'foo'
    with mock.patch.dict(os.environ, {constants.ENV_MODULE_NAME: module}):
        path = Helpers.search_path('default')
        with pytest.raises(Exception) as ex:
            find_func(path)
        assert f"Module {module} not found" in str(ex)

def test_function_not_found():
    func = 'foo'
    with mock.patch.dict(os.environ, {constants.ENV_FUNCTION_NAME: func}):
        path = Helpers.search_path('default')
        with pytest.raises(Exception) as ex:
            find_func(path)
        assert f"Function {func} not found" in str(ex)

# def test_find_func_invalid_handler():
#     os.environ["PYTHON_HANDLER"] = "handler.unknown"
#     handler_bar = Helpers.search_path("handler-handler")
#     with pytest.raises(Exception) as ex:
#         find_func(handler_bar)
#     assert "not found in module" in str(ex)

class Helpers:
    @staticmethod
    def search_path(case):
        return os.path.join(function_asset_base_dir, "http", case)
