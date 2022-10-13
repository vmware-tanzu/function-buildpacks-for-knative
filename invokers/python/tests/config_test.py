# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from pyfunc.config import Config

def test_init():
    cfg = Config('some/search/path', 'some_module', 'some_function')
    assert cfg.search_path == 'some/search/path'
    assert cfg.module_name == 'some_module'
    assert cfg.function_name == 'some_function'
