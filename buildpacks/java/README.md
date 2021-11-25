# Java Buildpack

This buildpack provides a way to convert your Java function into a container image.

## Getting started
To get started you'll need to create a directory where your function will be defined.

From within this directory we have to create a few files.
* WIP

## <a name="fp"></a> Accepted Function Parameters
The function handles either HTTP or CloudEvents based on the parameter's name and type. Only the following arguments are accepted:

(Note: This is not up-to-date and may be wrong.)

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
We've already created the builder for you: `us.gcr.io/daisy-284300/kn-fn/builder:0.0.3`

This builder can be used to create your function image. Firstly there are some tools you'll want

### Prerequisites
* [Buildpack CLI](https://buildpacks.io/docs/tools/pack/)

### <a name="usage"></a> Usage
Build the function container with the Buildpack CLI
```
pack build <your_image_name_and_tag> --builder us.gcr.io/daisy-284300/kn-fn/builder:0.0.3
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
If you want to quickly start writing your functions, take a look at the `templates/java` folder at the root of this repo.

## Examples
For some inspiration, take a look at the `samples/java` folder at the root of this repo.
