# Copyright 2021-2022 VMware, Inc.
# SPDX-License-Identifier: BSD-2-Clause

from fpdf import FPDF
from typing import Any
from urllib.parse import unquote_plus

import boto3
import os
import sys

def main(data: Any, attributes: dict):
    print(f"Got event data: {data}", file=sys.stderr)

    file = unquote_plus(data["s3"]["object"]["key"])
    bucket = data["s3"]["bucket"]["name"]
    if not file.endswith(".txt"):
        print(f"File {file} does not have .txt extension", file=sys.stderr)
        return attributes, "Can not convert non-txt file to PDF"

    print(f"Downloading {file} from bucket {bucket}", file=sys.stderr)
    response = client.get_object(Bucket=bucket, Key=file)

    pdf = FPDF()
    pdf.add_page()
    pdf.set_font("Arial", size = 15)
    for line in response['Body'].iter_lines():
        s = line.decode("utf-8")
        pdf.cell(200, 10, txt=s, ln=1, align='C')

    newName=file.replace(".txt", ".pdf")
    data = pdf.output(dest="S").encode("latin1") # Python > 3.0
    client.put_object(Body=data, Bucket=bucket, Key=newName)

    fresp = f"{file} in bucket {bucket} converted to PDF {newName}"
    print(fresp, file=sys.stderr)

    return attributes, fresp

class AWSCreds:
    def __init__(self) -> None:
        self._access_key = os.environ["AWS_ACCESS_KEY"]
        self._secret_key = os.environ["AWS_SECRET_KEY"]
        self._session_key = os.environ["AWS_SESSION_KEY"]
    
    @property
    def accessKey(self):
        return self._access_key
    
    @property
    def secretKey(self):
        return self._secret_key

    @property
    def sessionKey(self):
        return self._session_key

creds = AWSCreds()
client = boto3.client(
    's3',
    aws_access_key_id=creds.accessKey,
    aws_secret_access_key=creds.secretKey,
    aws_session_token=creds.sessionKey,
)
