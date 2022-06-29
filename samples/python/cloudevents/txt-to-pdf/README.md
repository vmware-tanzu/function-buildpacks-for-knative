# txt-to-pdf

## Summary

This Python function listens for AWS S3 create events. When a new `.txt` file is detected, it will attempt to convert the text file into a PDF, then upload it to S3.

## Prerequisites
* Cloud Native Runtimes: 1.0.0+
* AWS CLI
* Docker
* Buildpack CLI
* Kapp
* Ytt

## Known Issues
* n/a

## Demo
1. Create a bucket in S3 for the demo. Note down the ARN of this S3 bucket.

1. Get your AWS accesskey and secretkey that will be used by Triggermesh to listen for events.

1. We will also be using STS to [generate a temporary session](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html) for the application via the AWS CLI. This will generate the three keys you use under app in `creds.yaml` below.
    ```
    aws sts get-session-token --duration-seconds=129600 # This will generate a session token for the maximum allowed time of 1.5 days
    ```

1. Now let's create a new file called `creds.yaml` in the current folder (`samples/python/cloudevent/txt-to-pdf`):
    ```
    cat > creds.yaml << EOF
    ---
    triggermesh:
      accesskey: <your access key from step 2>
      secretkey: <your secret key from step 2>

    app:
      accesskey: <your access key from step 3 (This value is different from the one in step 2!)>
      secretkey: <your secret key from step 3 (This value is different from the one in step 2!)>
      sessionkey: <your session key from step 3>

    bucket_arn: <your bucket ARN from step 1 (Must be the full ARN including the region and account)>
    EOF
    ```

1.  If you want to change the location of the function, you need to define the environment variable `FUNCTION_IMAGE`
    ```
    export FUNCTION_IMAGE=<your full image url that can be pushed to>
    ```

1. Create the function container image.
    ```
    make build
    ```

1. Publish your image.
    ```
    make publish
    ```

1. Deploy your app!
    ```
    make deploy
    ```

### Cleanup
To cleanup, simply run:
```
make destroy
```    

## Other

[Video demo link](https://vmware.enterprise.slack.com/files/WS0819VJM/F02C0ASMJAY/func-demo-full.mp4?origin_team=T024JFTN4&origin_channel=C021B90DLMA)
