# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from typing import Any
from urllib.parse import unquote_plus
from cryptography.fernet import Fernet

import binascii
import boto3
import os
import sys
import json

def DoEvent(data: Any, attributes: dict):

    body = json.loads(data["Body"])
    
    file = s3Client.get_object(Bucket=body["bucket"], Key=body["key"])
    fernet = Fernet(config.encryption_key)
    
    full_body = bytearray()
    for bytes_chunks in file['Body'].iter_chunks():
        full_body += bytes_chunks    
    encrypted_body = fernet.encrypt(binascii.hexlify(full_body))
    s3Client.put_object(Body=encrypted_body, Bucket=config.bucket_name, Key=body["key"])
    
    s3Client.delete_object(Bucket=body["bucket"], Key=body["key"])

    fresp = f"{body['key']} in bucket {body['bucket']} encrypted"
    return attributes, fresp

class Config:
    def __init__(self) -> None:
        self._access_key = os.environ["AWS_ACCESS_KEY"]
        self._secret_key = os.environ["AWS_SECRET_KEY"]
        self._aws_s3_bucket = os.environ["AWS_S3_BUCKET_NAME"]
        self._encryption_key = Fernet.generate_key()

    @property
    def access_key(self):
        return self._access_key

    @property
    def secret_key(self):
        return self._secret_key

    @property
    def bucket_name(self):
        return self._aws_s3_bucket

    @property
    def encryption_key(self):
        return self._encryption_key


config = Config()
s3Client = boto3.client(
    's3',
    aws_access_key_id=config.access_key,
    aws_secret_access_key=config.secret_key,
)
