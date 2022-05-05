# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

FROM python@sha256:301797d06f5828195f496f1f3022e370d5743e44044e56656f45c4a3c4174ca4

RUN python -m pip install --upgrade pip
RUN pip install tox pytest

COPY . /workspace/invoker
WORKDIR /workspace/invoker
ENV PYTHONPATH "${PYTHONPATH}:/workspace/invoker"
RUN tox sdist
RUN mkdir -p /out
RUN cp /workspace/invoker/.tox/dist/*.tar.gz /out

ENTRYPOINT [ "tox" ]
CMD [ "tests" ]

# We're currently generating the sha from outside the docker container.
# WORKDIR /out
# RUN find . -type f -iname '*.tar.gz' -exec sh -c 'shasum -a 256 {} > {}.sha256' \;
