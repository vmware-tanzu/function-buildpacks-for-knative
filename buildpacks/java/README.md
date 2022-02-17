# Java Function Buildpack

The Java Function Buildpack is a Cloud Native Buildpack that provides a Spring Boot application for executing functions.

## Behaviour
This buildpack will participate if any of the following conditions are met:
* A file with the name `func.yaml` is detected

The buildpack will do the following if detection passed:
* Request for a JRE to be installed
* Contributes the Spring Boot application to a layer marked `launch` with the layer's path prepended to `$CLASSPATH`
* Contributes environment variables defined in `func.yaml` to the `launch` layer

## Getting started
To get started you'll need to create a directory where your function will be defined.

From within this directory we require a few files to properly detect this as a Java function:
* `func.yaml`: We use this to configure the runtime environment variables. See the [Knative Func CLI docs](https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/guides/func_yaml.md) for more details.
* `pom.xml` or `build.gradle`: These are used by the other Java buildpacks to compile your function.
* Java package in folder `src/main/java/functions`: This is the default location your function will be detected. If you do choose to use another package to store your functions, you will need to [set a new search location](#TODO).

## Compiling Your Function
To compile your function with the buildpack, we've provided a builder which has all the pre-requisites ready to go.
You can find it [on github](https://github.com/vmware-tanzu/function-buildpacks-for-knative/pkgs/container/function-buildpacks-for-knative%2Ffunctions-builder).

```
ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder
```

### Prerequisites
* [Buildpack CLI](https://buildpacks.io/docs/tools/pack/)

### <a name="usage"></a> Usage
Build the function container with the Buildpack CLI
```
pack build <your_image_name_and_tag> --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:<version>
```

Publish it to your registry:
```
docker push <your_image_name_and_tag>
```

## Templates
If you want to quickly start writing your functions, take a look at the `templates/java` folder at the root of this repo.

## Examples
For some inspiration, take a look at the `samples/java` folder at the root of this repo.
