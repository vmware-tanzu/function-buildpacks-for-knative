# Python Function Buildpack

The Python Function Buildpack is a Cloud Native Buildpack that provides a Python function invoker application for executing functions.

## Behaviour
This buildpack will participate if any of the following conditions are met:
- The `BP_FUNCTION` environment variable is set.
- A valid `func.yaml` exists in the function directory. See [func.yaml](#func.yaml) 

The buildpack will do the following if detection passed:
* Request for a Python runtime to be installed to a layer marked `build` and `launch`
* Request for a `pip` to be installed to a layer marked `build`
* Contributes the Python function invoker application to a layer marked `launch`
* Contributes environment variables defined in `func.yaml` to the `launch` layer
* Contributes a validation layer which is used to determine if the function is properly defined

## Getting started
To get started you'll need to create a directory where your function will be defined.

From within this directory we require a few files to properly detect this as a Python function:
* `func.py`: This python module will be where we search for a function by default.
  * If you want to use a different name for the file. See [configuration](#configuration) or [`func.yaml`](#func.yaml).
  * This file should contain the function to invoke when we receive an event.
    * The function can handle http requests:
      ```
      from typing import Any

      def main(req: Any):
        return "Handled request!"
      ```
    * The function can handle CloudEvents:
      ```
      from typing import Any

      def main(data: Any, attributes: dict):
        return attributes, "Handled cloudevent!"
      ```
    * You can find more details about the different accepted parameters [below](#fp).

* <a name="func.yaml"></a>`func.yaml` (optional): This is the configuration used to configure your function.

  The python module and function name can be modified here by defining some environment variables in the `envs` section.
  ```
  name: test
  runtime: python
  envs:
  - name: MODULE_NAME
    value: my_module
  - name: FUNCTION_NAME
    value: my_func
  ```
  By defining the above, we will look for a `my_module.py` instead of `func.py` which contains a function with the name `my_func`.
 
  This buildpack makes use of `envs` and `options`. The keys `name` and `runtime` are required to maintain compatibility with Knative func cli, but are not used by this buildpack. 
  See [Knative's func.yaml documentation](https://github.com/knative/func/blob/main/docs/reference/func_yaml.md) 
  for more `func.yaml` information.

  **NOTE**: The environment variables here (namely `MODULE_NAME` and `FUNCTION_NAME` will be overriden by the values specified by `BP_FUNCTION`)


* `requirements.txt`: This file is required by the Python dependency. It is used to define your function's dependencies. If you do not have any, you still need to provide an empty file.

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
To compile your function with the buildpack, we've provided a builder which has all the pre-requisites ready to go.
You can find it [on github](https://github.com/vmware-tanzu/function-buildpacks-for-knative/pkgs/container/function-buildpacks-for-knative%2Ffunctions-builder).

```
ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder
```
### Prerequisites
* [Buildpack CLI](https://buildpacks.io/docs/tools/pack/)

### <a name="usage"></a> Usage

## <a name="configuration"></a> Configuration

| Environment Variable | Description |
| -------------------- | ----------- |
| `$BP_FUNCTION` | Configure the function handler.  Defaults to `func.main`. |

Build the function container with the Buildpack CLI
```
pack build <your_image_name_and_tag> --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:<version>
```

Publish it to your registry:
```
docker push <your_image_name_and_tag>
```

## Health Endpoints

The Python Invoker has health endpoints exposed by `healthz`. By default, the path is found at `localhost:8080/healthz/live` or `localhost:8080/healthz/ready`.

## Templates
If you want to quickly start writing your functions, take a look at the `templates/python` folder at the root of this repo.

## Examples
For some inspiration, take a look at the `samples/python` folder at the root of this repo.

CloudEvent samples:
- [Simple S3 Interaction](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/cloudevent/s3-lambda)
- [SQS Secrets Encrypter](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/cloudevent/sqs-lambda)
- [Automatic S3 txt-to-pdf Converter](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/cloudevent/txt-to-pdf)

HTTP samples:
- [Link](https://gitlab.eng.vmware.com/daisy/functions/buildpacks/-/tree/master/samples/python/http)
