# txt-to-pdf

This is the core txt-to-pdf sample folder.

## Summary

This Java function listens for AWS S3 create events. When a new `.txt` file is detected, it will attempt to convert the text file into a PDF, then upload it to S3.

## Video Demos

- [Python video demo](https://vmware.enterprise.slack.com/files/WS0819VJM/F02C0ASMJAY/func-demo-full.mp4?origin_team=T024JFTN4&origin_channel=C021B90DLMA)
- [Java video demo](#)

## Setup

### Prerequisites
* Cloud Native Runtimes: 1.0.0+
* AWS CLI
* Docker
* Buildpack CLI
* Kapp
* ytt >= 0.36

### Known Issues
* n/a

### Instructions

1. Create a bucket in S3 for the demo. Note down the ARN of this S3 bucket.

1. Get your AWS access key and secret key that will be used by Triggermesh to listen for events.

1. Create a new file called `creds.yaml` in the current folder (`samples/java/cloudevent/txt-to-pdf`):
    ```
    cat > creds.yaml << EOF
    ---
    triggermesh:
      accesskey: <your access key from step 2>
      secretkey: <your secret key from step 2>

    app:
      accesskey: <your access key from step 2>
      secretkey: <your secret key from step 2>

    bucket_arn: <your bucket ARN from step 1 (Must be the full ARN including the region and account)>
    EOF
    ```

    Optionally, if you need to use STS to configure your AWS credentials, see "AWS Credentials via STS" below.

1.  If you want to change the location of the function, you need to define the environment variable `FUNCTION_IMAGE`
    ```
    export FUNCTION_IMAGE=<your full image url that can be pushed to>
    ```

1. (Optional) Deploy Cloud Native Runtimes on a platform of your choice, including [these verifying steps](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.0/tanzu-cloud-native-runtimes-1-0/GUID-verifying-triggermesh.html).

1. Create a TriggerMesh S3 Source mapped to the public bucket name you created in Step 1.

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

If you encounter any Knative errors while re-deploying the app, be sure to delete the `ksvc` consumer before re-running `make deploy`.


## AWS Credentials via STS

1. We will also be using STS to [generate a temporary session](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html) for the application via the AWS cli.
    ```
    aws sts get-session-token --duration-seconds=1800 # This will generate a session that is going to last for 30m
    ```

1. Now let's create a new file called `creds.yaml` in the current folder (`samples/java/cloudevent/txt-to-pdf`):
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

## Legacy Deployment YAML

apiVersion: v1
kind: Secret
metadata:
  name: my-aws-creds
  namespace: cnr-demo
type: Opaque
stringData:
  aws-access-key: "key"
  aws-secret-key: "key"
  aws-session-key: "key"
---
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: consumer
  namespace: cnr-demo
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/minScale: "1"
    spec:
      containers:
        - image: <your-image>
          imagePullPolicy: Always
          env:
            - name: FLASK_ENV
              value: development
            - name: AWS_ACCESS_KEY_ID
              valueFrom:
                secretKeyRef:
                  name: my-aws-creds
                  key: aws-access-key
            - name: AWS_SECRET_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  name: my-aws-creds
                  key: aws-secret-key
            - name: AWS_SESSION_TOKEN
              valueFrom:
                secretKeyRef:
                  name: my-aws-creds
                  key: aws-session-key
```
