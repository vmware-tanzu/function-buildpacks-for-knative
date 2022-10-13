# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause


class Config:
    def __init__(self, search_path: str, module_name: str, function_name: str) -> None:
        self._module_name = module_name
        self._function_name = function_name
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
