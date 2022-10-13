# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

import os
import typing

import pytest
from pyfunc import find_func

function_asset_base_dir = os.path.join("tests", "functions")

def test_find_func(get_search_path):
    path = get_search_path('locate')
    func = find_func(path, 'some_module', 'some_function')
    assert isinstance(func, typing.Callable)

def test_module_not_found(get_search_path):
    path = get_search_path('locate')
    with pytest.raises(Exception, match=f"Module 'not_a_module' not found in '{path}'"):
        find_func(path, 'not_a_module', '')

def test_function_not_found(get_search_path):
    path = get_search_path('locate')
    with pytest.raises(Exception, match=f"Function 'not_a_function' not found in module 'some_module'"):
        find_func(path, 'some_module', 'not_a_function')

@pytest.fixture
def get_search_path():
    def fn(case: str):
        return os.path.join(function_asset_base_dir, "http", case)
    return fn
