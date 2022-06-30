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

## Setup

1. Create a public S3 bucket for the demo, and note the ARN of this S3 bucket -- you'll use it soon. (e.g. `arn:aws:s3:us-west-2:123456789012:bucket_name`)

1. Deploy a Kubernetes cluster (recommended: with Tanzu Application Platform installed)

    -  If you are NOT using Tanzu Application Platform, install [Cloud Native Runtimes](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.2/tanzu-cloud-native-runtimes/GUID-install.html) first.

    -  If you ARE using Tanzu Application Platform, you will need to configure the Tanzu Image Registry secrets (Steps WIP).
    
1. [Verify](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.2/tanzu-cloud-native-runtimes/GUID-verify-installation.html) your Cloud Native Runtimes installation was successful. 

    >  Important: Ensure you verify [TriggerMesh SAWS](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.2/tanzu-cloud-native-runtimes/GUID-verifying-triggermesh.html) up to Step 5, then instead proceed with deplying an AWSS3Source.

1. Deploy an [AWS S3 Source](https://github.com/triggermesh/aws-event-sources/blob/main/config/samples/awss3source.yaml)

    - (Optional) Learn more about [Sources and TriggerMesh](https://docs.triggermesh.io/cloud/sources/awss3/), since Cloud Native Runtimes uses TriggerMesh SAWS to subscribe to event notifications and connect to AWS services.

    - You may use the template below, which already contains `namespace` added. Be sure to replace `<YOUR-ARN`> with the ARN from Step 1.
        ```
        kubectl apply -f - << EOF
        apiVersion: sources.triggermesh.io/v1alpha1
        kind: AWSS3Source
        metadata:
        name: sample
        spec:
        arn: <YOUR-ARN>

        eventTypes:
        - s3:ObjectCreated:*
        - s3:ObjectRemoved:*

        credentials:
            accessKeyID:
            valueFromSecret:
                name: awscreds
                key: aws_access_key_id
            secretAccessKey:
            valueFromSecret:
                name: awscreds
                key: aws_secret_access_key

        sink:
            ref:
            apiVersion: eventing.knative.dev/v1
            kind: Broker
            name: default
            namespace: ${WORKLOAD_NAMESPACE}
        EOF
        ```

1. Obtain your AWS `accesskey` and `secretkey` that will be used in the next steps.

1. Use the following STS command to to [generate a temporary session](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_use-resources.html) for the application via the AWS CLI. This will generate three keys you use under app in `creds.yaml` below, please note them down somewhere safe temporarily.
    ```
    aws sts get-session-token --duration-seconds=129600 # This will generate a session that is going to last for 36h
    ```

1. Now let's create a new file called `creds.yaml` in the current folder (`samples/python/cloudevent/txt-to-pdf`). Paste the following into your terminal to create this file, then fill in your keys.
    ```
    cat > creds.yaml << EOF
    ---
    triggermesh:
        accesskey: <your AccessKeyId from step 2>
        secretkey: <your SecretAccessKey from step 2>

    app:
        accesskey: <your AWS_ACCESS_KEY_ID from step 3 (different from the value in step 2!)>
        secretkey: <your AWS_ACCESS_KEY_ID from step 3 (different from the value in step 2!)>
        sessionkey: <your SessionToken from step 2>

    bucket_arn: <your bucket ARN from step 1 (Must be the full ARN including the region and account)>
    EOF
    ```

1.  (Optional) If you want to change the location of the function, you need to define the environment variable `FUNCTION_IMAGE`
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
