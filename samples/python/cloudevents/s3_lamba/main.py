# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from typing import Any
from urllib.parse import unquote_plus

import boto3
import os
import sys

def DoEvent(data: Any, attributes: dict):

    bucket = data["s3"]["bucket"]["name"]
    key = unquote_plus(data["s3"]["object"]["key"], encoding='utf-8')
    tagging= "STATE=DELETED&FROM={}".format(bucket)
    more_binary_data = b'SOMEONE DELETED ME! :('
    client.put_object(
        Bucket=bucket,
        Key=key,
        Body=more_binary_data,
        Tagging=tagging
    )

    return key + " delete from " + bucket

class AWSCreds:
    def __init__(self) -> None:
        self._access_key = os.environ["AWS_ACCESS_KEY"]
        self._secret_key = os.environ["AWS_SECRET_KEY"]

    @property
    def accessKey(self):
        return self._access_key

    @property
    def secretKey(self):
        return self._secret_key


creds = AWSCreds()
client = boto3.client(
    's3',
    aws_access_key_id=creds.accessKey,
    aws_secret_access_key=creds.secretKey)
