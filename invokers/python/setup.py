# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from os import path
from setuptools import find_packages, setup

basepath = path.dirname(__file__)
desc_file = path.abspath(path.join(basepath, "README.md"))
version_file = path.abspath(path.join(basepath, "VERSION"))
modules = [
    path.abspath(path.join(basepath, "pyfunc"))
]

with open(desc_file, 'r') as fh:
    long_description = fh.read()

with open(version_file, 'r') as fh:
    version = fh.read()

setup(
    name='pyfunc-invoker',
    packages=find_packages(include=['pyfunc']),
    version=version,
    long_description=long_description,
    long_description_content_type='text/markdown',
    license='Apache',
    install_requires=[
        'cloudevents >=1.2, <2',
        'flask-healthz >=0.0.3',
        'Flask >=2,<3',
        'waitress >=2',
    ],
    setup_requires=['pytest-runner'],
    tests_require=['pytest==4.4.1'],
    test_suite='tests',
    classifiers=[
        "Programming Language :: Python :: 3",
        'License :: OSI Approved :: Apache Software License',
        'Operating System :: Unix',
    ],
)
