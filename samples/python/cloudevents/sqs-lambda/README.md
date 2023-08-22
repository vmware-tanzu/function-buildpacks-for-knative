# Secret Encrypter Example

We want to demonstrate how a chain of functions can live in different clouds.

In this example, a fake organization needs to keep secret files encrypted in the cloud, by using a secret key they feel uncomfortable sharing with cloud providers.

When a new secret file (for this example, will be identified when the file name includes the word `secret`) is uploaded, it will send a message to the private encrypter that is our Knative function instance. This encrypter will then: get, encrypt, and upload the encrypted file, then delete the unencrypted version of it.


## Architecture

![Architecture Flow](assets/arch.jpeg?raw=true "Encrypter Data Plane")

## Pre-requisites

1. Have an AWS account where you can:
   1. Add new users and get `account access id` and `secret access key`.
   1. Create and deploy Lambda functions
   1. Create or have access to SQS instances
1. A cluster with CNR installed
1. Docker installed locally

## Step-by-step

### Create your S3 bucket (entry point)

This first bucket `Item Storage` is the bucket you will use to trigger the whole flow, as it is the only manual entry point for this demo.

### Create your SQS instance

Create a queue where the AWS Lambda function can send messages and notify the system something needs attention.

### Create your S3 secrets instance

Create another S3 bucket where the Knative function will storage the encrypted files.

### Deploy your AWS Lambda function

This Lambda function will read objects from the `Item Storage` and publish to a queue that a new object has arrived.

This Lambda will filter objects that has the word `secret` in their name and push to the queue that the encrypter is listening.

Code for your lambda is on `./s3-trigger-to-sqs-lambda`. (I hope these name make sense to you :) - dolfo)

#### Steps

// TODO We might clean the policies and only have write for SQS and execution for Lambda

1. Run `npm install` in the folder `s3-trigger-to-sqs-lambda`.
1. Zip the files of the function folder and upload it to lambda functions. `zip -r NAME.zip path/of/lambda/folder`
1. Add environment variables to the lambda function code in configuration:
   - `QUEUE_URL` the url of the SQS we created.
   - `REGION` the region where the SQS instance is deployed.
1. Under `configurations -> permissions` tab we need to add a new role that has `AmazonSQSFullAccess` and `AWSLambda_FullAccess` policies attached.

_Note: This is in NodeJS but you can deploy in any language you like. The important part are the permissions._

### Deploy your Knative function

#### Docker Steps (Required)

The folder `knative-function-encrypter` has the function that will be listening to our SQS queue.

With docker running, build the image from this folder:

```
pack build encrypter --path PATH/TO/knative-function-encrypter --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.4.3 --env BP_FUNCTION=main.main
```

Tag and push to your registry:
```
docker tag encrypter REGISTRY.REPO/URL/FOR/IMAGE && docker push REGISTRY.REPO/URL/FOR/IMAGE 
```

After you have the image, you can deploy your Knative Serving with the secret of your AWS role that has read and push access to SQS and S3.

#### Knative Steps (Required)

First, create a namespace in your Kubernetes cluster and export an environment variable for convenience:
```
export WORKLOAD_NAMESPACE='sqs-lambda-demo'
kubectl create namespace ${WORKLOAD_NAMESPACE}
```

Save a copy of the below `.yaml` files to your local machine, then edit the # CHANGE ME lines. Afterwards, `k apply -f MYFILE.yaml` to apply them to your cluster.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-aws-creds
  namespace: ${WORKLOAD_NAMESPACE}
type: Opaque
stringData:
  aws-access-key-id: "[ACCESS KEY ID]" # CHANGE ME
  aws-secret-access-key: "[SECRET ACCESS KEY]" # CHANGE ME
---
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: encrypter
  namespace: ${WORKLOAD_NAMESPACE}
spec:
  template:
    metadata:
      annotations:
        autoscaling.knative.dev/minScale: "1"
    spec:
      containers:
        - image: REGISTRY.REPO/URL/FOR/IMAGE # CHANGE ME
          imagePullPolicy: Always
          env:
            - name: FLASK_ENV
              value: development
            - name: AWS_ACCESS_KEY
              valueFrom:
                secretKeyRef:
                  name: my-aws-creds
                  key: aws-access-key-id
            - name: AWS_SECRET_KEY
              valueFrom:
                secretKeyRef:
                  name: my-aws-creds
                  key: aws-secret-access-key
            - name: AWS_S3_BUCKET_NAME
              value: [BUCKET NAME FOR ENCRYPTED FILES] # CHANGE ME
```

Our function is ready to get events! We need to deploy the source and trigger of our function:

```yaml
apiVersion: eventing.knative.dev/v1
kind: Broker
metadata:
  name: default
  namespace: ${WORKLOAD_NAMESPACE}
---
apiVersion: sources.triggermesh.io/v1alpha1
kind: AWSSQSSource
metadata:
  name: sqs-invoice-source
  namespace: ${WORKLOAD_NAMESPACE}
spec:
  arn: "arn:aws:sqs:[REGION]:[ACCOUNT_ID]:[SQS INSTANCE NAME]" # CHANGE ME
  credentials:
    accessKeyID:
      valueFromSecret:
        name: my-aws-creds
        key: aws-access-key-id
    secretAccessKey:
      valueFromSecret:
        name: my-aws-creds
        key: aws-secret-access-key

  sink:
    ref:
      apiVersion: eventing.knative.dev/v1
      kind: Broker
      name: default
      namespace: ${WORKLOAD_NAMESPACE}

---
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: trigger
  namespace: ${WORKLOAD_NAMESPACE}
spec:
  broker: default
  subscriber:
    ref:
      apiVersion: v1
      kind: Service
      name: encrypter
      namespace: ${WORKLOAD_NAMESPACE}
```

If we did everything correctly, now you can go to S3 and upload a file with the name `secret` in their file name and watch.
