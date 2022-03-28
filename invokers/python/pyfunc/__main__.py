# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from pyfunc import invoke, find_func

import sys
from argparse import ArgumentParser


def check(args):
    try:
      find_func(args.search_path)
    except Exception as e:
      print("Failed to detect function: " + str(e), file=sys.stderr)
      sys.exit(1)

def start(args):
    invoke.main(args.search_path)

parser = ArgumentParser(prog='pyfunc')
parser.set_defaults(func=lambda args: parser.print_help())
subparsers = parser.add_subparsers(help='sub-command help')

parser_check = subparsers.add_parser('check', help='check if the module and function can be loaded')
parser_check.add_argument('-s', '--search_path', type=str, default='.')
parser_check.set_defaults(func=check)

parser_start = subparsers.add_parser('start', help='start the python invoker')
parser_start.add_argument('-s', '--search_path', type=str, default='.')
parser_start.set_defaults(func=start)

args = parser.parse_args(sys.argv[1:])
args.func(args)
