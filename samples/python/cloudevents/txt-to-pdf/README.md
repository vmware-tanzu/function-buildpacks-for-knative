# txt-to-pdf

## Summary

This Python function listens for AWS S3 create events. When a new `.txt` file is detected, it will attempt to convert the text file into a PDF, then upload it to S3.

## Video Demos

- [Python video demo](https://vimeo.com/724580619)
- [Java video demo](https://vimeo.com/724580576)

## Prerequisites
* [kubectl](https://kubernetes.io/docs/tasks/tools/)
* [aws](https://aws.amazon.com/cli/)
* [docker](https://docs.docker.com/engine/install/)
* [pack](https://buildpacks.io/docs/tools/pack/)
* [kapp](https://carvel.dev/kapp/)
* [ytt](https://carvel.dev/ytt/)

## Known Issues
* Image registry authentication syncing with Tanzu Application Platform
* ⚠️ Cloud Native Runtimes has a regression in which TriggerMesh does not work

## Setup

1. Log in to AWS console and [obtain your AWS Access Key and Secret Key](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html), save these somewhere secure.

1. [Create a public S3 bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/creating-bucket.html) for the demo, and note the [ARN](https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html) of this S3 bucket.

1. Deploy a Kubernetes cluster with [Tanzu Application Platform (TAP)](https://docs.vmware.com/en/VMware-Tanzu-Application-Platform/1.1/tap/GUID-install-intro.html) installed
    
    -  Ensure you have [Tanzu Image Registry](https://docs.vmware.com/en/VMware-Tanzu-Application-Platform/1.1/tap/GUID-install.html) secrets configured properly

    -  Install [Cloud Native Runtimes](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.2/tanzu-cloud-native-runtimes/GUID-install.html)
    
1. [Verify](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.2/tanzu-cloud-native-runtimes/GUID-verify-installation.html) your Cloud Native Runtimes installation was successful. 

    >  Important: Ensure you verify [TriggerMesh SAWS](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.2/tanzu-cloud-native-runtimes/GUID-verifying-triggermesh.html), as any error in the cluster or configuration should surface by here.

1. Locate the `values.yaml` file under this sample's `/config`. You will use your AWS Access Key and AWS Secret Key to fill out the fields here. Use the following below as a guide:
    ```
    #@data/values
    ---
    namespace: cnr-demo

    triggermesh:
        accesskey: <your AWS_ACCESS_KEY_ID>
        secretkey: <your AWS_SECRET_KEY_ID>

    app:
        accesskey: <your AWS_ACCESS_KEY_ID>
        secretkey: <your AWS_SECRET_KEY_ID>

    function_image: <optional>
    bucket_arn: <e.g. arn:aws:s3:us-west-2:123456789012:bucket_name>
    ```

1.  (Optional) Change the function_image repository location

## (Outdated / Optional) Using AWS STS for Temporary Sessions

1. Obtain your AWS `accesskey` and `secretkey` that will be used in the next steps.

1. Use the following STS command to to [generate a temporary session](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html) for the application via the AWS CLI. This will generate three keys you use under app in `creds.yaml` below, please note them down somewhere safe temporarily.
    ```
    aws sts get-session-token --duration-seconds=129600 # This will generate a session that is going to last for 36h
    ```

1. Use the following template as a reference to fill out `config/values.yaml`:
    ```
    triggermesh:
        accesskey: <your AWS Access Key>
        secretkey: <your AWS Secret Key>

    app:
        accesskey: <your STS access key>
        secretkey: <your STS secret key>
        sessionkey: <your STS session key>

    bucket_arn: <your full bucket ARN>
    ```

1. (Optional) If you want to change the location of the function, you need to define the environment variable `FUNCTION_IMAGE`
    ```
    export FUNCTION_IMAGE=<your full image url that can be pushed to>
    ```

1. Uncomment the session key field found under `config/300-app-secret.yaml`
## Deploying

There are three steps in our Makefile to build and deploy the demo. For your convenience:

```
make build && make publish && make deploy
```

## Demo'ing

1. Upload a file ending in `.txt` to your public S3 bucket. (We've provided you `story.txt`)

2. A `.pdf` should have generated in the same bucket immediately after the upload.

3. Download and view your newly converted document!
 

### Cleanup
To cleanup, simply run:
```
make destroy
```    

If you encounter any Knative errors while re-deploying the app, be sure to delete the `ksvc` consumer before re-running `make deploy`.
