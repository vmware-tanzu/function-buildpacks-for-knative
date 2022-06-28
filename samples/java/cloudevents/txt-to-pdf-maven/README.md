# txt-to-pdf

## Summary

This Java function listens for AWS S3 create events. When a new `.txt` file is detected, it will attempt to convert the text file into a PDF, then upload it to S3.

## Video Demos

- [Python video demo](https://vimeo.com/724580619)
- [Java video demo](https://vimeo.com/724580576)

## Prerequisites
* Tanzu Application Platform
    * Cloud Native Runtimes installed
* AWS CLI
* Docker
* Buildpack CLI
* Kapp
* Ytt

## Known Issues
* Image registry authentication syncing with Tanzu Application Platform

## Setup

1. Create a public S3 bucket for the demo, and note the ARN of this S3 bucket -- you'll use it soon. (e.g. `arn:aws:s3:us-west-2:123456789012:bucket_name`)

1. Deploy a Kubernetes cluster (recommended: with Tanzu Application Platform)

    1. Install [Cloud Native Runtimes](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.2/tanzu-cloud-native-runtimes/GUID-install.html), up to Step 4.

    1. Deploy an [AWS S3 Source](https://github.com/triggermesh/aws-event-sources/blob/main/config/samples/awss3source.yaml) instead of a CodeCommit source, as in Step 5.

    1. Modify your S3 Source to include `${WORKLOAD_NAMESPACE}` as needed in the Cloud Native Runtimes docs and the S3 ARN information.

1. Obtain your AWS `accesskey` and `secretkey` that will be used in the next steps.

1. Use the following STS command to to [generate a temporary session](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html) for the application via the AWS CLI. This will generate the three keys you use under app in `creds.yaml` below.
    ```
    aws sts get-session-token --duration-seconds=1800 # This will generate a session that is going to last for 30m
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

## Deploying

There are three steps in our Makefile to build and deploy the demo. For your convenience:

```
make build && make publish && make deploy
```

## Demo'ing

1. Upload a file ending in `.txt` to your public S3 bucket.

2. A `.pdf` should have generated in the same bucket immediately after the upload.

3. Download and view your newly converted document!
 

### Cleanup
To cleanup, simply run:
```
make destroy
```    
