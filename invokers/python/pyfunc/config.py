# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from .constants import *
import os

class Config:
    def __init__(self, search_path=SEARCH_PATH_DEFAULT) -> None:
        self._module_name = os.getenv(ENV_MODULE_NAME, MODULE_NAME_DEFAULT)
        self._function_name = os.getenv(ENV_FUNCTION_NAME, FUNCTION_NAME_DEFAULT)
        self._search_path = search_path

    def __repr__(self) -> str:
        return 'Config:\n' \
            f'  Function {self.function_name}\n' \
            f'  Module {self.module_name}\n' \
            f'  Search Path {self.search_path}'

    @property
    def module_name(self) -> str:
        return self._module_name

    @property
    def function_name(self) -> str:
        return self._function_name

    @property
    def search_path(self) -> str:
        return self._search_path
