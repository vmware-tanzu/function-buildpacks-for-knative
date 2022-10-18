# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

FROM maven:3.8.3-jdk-11

WORKDIR /usr/src/invoker
COPY . /usr/src/invoker

RUN ./package.sh
RUN mkdir -p /out && cp target/java-function-invoker-$(cat ./VERSION).jar /out

ENTRYPOINT [ "mvn" ]
CMD [ "test" ]
