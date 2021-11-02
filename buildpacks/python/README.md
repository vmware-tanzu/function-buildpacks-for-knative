# Python Buildpack

This buildpack provides a way to convert your python3 function into a container image.

## Getting started
To get started you'll need to create a directory where your function will be defined.

From within this directory we have to create a few files.
* <a name="project.toml"></a>`project.toml`: This is the configuration used to give the buildpack some configurations.
  * The python module and function name can be modified here by defining a new build environment variable.
    ```
    [[build.env]]
    name="PYTHON_HANDLER"
    value="main.DoEvent"
    ```
    By defining the above instead of a `handler.py` file like below, the file should now be `main.py` containing a function with the name `DoEvent`

* `handler.py`: This python module will be where we search for a function by default.
  * If you want to use a different name for the file. See description for `project.toml`.
  * This file should contain the function to invoke when we receive an event.
    * The function can handle http requests:
      ```
      from typing import Any

      def handler(req: Any):
        return "Handled request!"
      ```
    * The function can handle CloudEvents:
      ```
      from typing import Any

      def handler(data: Any, attributes: dict):
        return attributes, "Handled cloudevent!"
      ```
    * You can find more details about the different accepted parameters [below](#fp).

* `requirements.txt`: This file is used for defining your dependencies. However if you have no dependencies, we're still expecting an empty file.
  * TODO: Remove the expectation of file `requirements.txt`

## <a name="fp"></a> Accepted Function Parameters
The function handles either HTTP or CloudEvents based on the parameter's name and type. Only the following arguments are accepted:
| name | request type | description | details |
|-|-|-|-|
| event | CloudEvent | Entire CloudEvent object | event |
| data | CloudEvent | Data portion of CloudEvent object | event.data |
| payload | CloudEvent | Data portion of CloudEvent object | event.data |
| attributes | CloudEvent | All CloudEvent keys and values as dictionary | |
| req | HTTP | Entire HTTP request (flask) | request |
| request | HTTP | Entire HTTP request (flask) | request |
| body | HTTP | Body of HTTP request (flask) | request.get_data() |
| headers | HTTP | HTTP request (flask) headers | request.headers |

## Compiling Your Function
We've already created the builder for you: `us.gcr.io/daisy-284300/kn-fn/builder:0.0.1`

This builder can be used to create your function image. Firstly there are some tools you'll want

### Prerequisites
* [Buildpack CLI](https://buildpacks.io/docs/tools/pack/)

### <a name="usage"></a> Usage
Build the function container with the buildpack cli
```
pack build <your_image_name_and_tag> --builder us.gcr.io/daisy-284300/kn-fn/builder:0.0.1
```

Publish it to your registry:
```
docker push <your_image_name_and_tag>
```

Deploy it to your cluster!
* If you're using Knative make sure to setup your eventing triggers and set
  1. Create your Knative service:
      ```
      apiVersion: serving.knative.dev/v1
      kind: Service
      metadata:
        name: consumer
      spec:
        template:
          spec:
            containers:
              - image: <your_image_name_and_tag>
      ```
  1. Follow the instructions on [Knative's eventing documentation](https://knative.dev/docs/eventing/broker/) about targeting your consumer service
* If you're deploying just an HTTP function then you can deploy it via a [deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) with the appropriate [service](https://kubernetes.io/docs/concepts/services-networking/service/).

## Templates
If you want to quickly start writing your functions, take a look at the `templates/python` folder at the root of this repo.

## Examples
For some inspiration, take a look at the `samples/python` folder at the root of this repo.

CloudEvent samples:
- [Simple S3 Interaction](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/cloudevent/s3_lamba)
- [SQS Secrets Encrypter](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/cloudevent/sqs-lambda)
- [Automatic S3 txt-to-pdf Converter](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/cloudevent/txt-to-pdf)

HTTP samples:
- [Link](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/http)
