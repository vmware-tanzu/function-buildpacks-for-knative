# Java Function Buildpack

The Java Function Buildpack is a Cloud Native Buildpack that provides a Spring Boot application for executing functions.

## Behaviour
This buildpack will participate if any of the following conditions are met:
* A buildpack configuration variable `BP_FUNCTION` is explicitly set.
* A file with the name `func.yaml` is detected.

The buildpack will do the following if detection passed:
* Request for a JRE to be installed
* Contributes the function invoker to a layer marked `launch` with the layer's path prepended to `$CLASSPATH`
* Contributes environment variables defined in `func.yaml` to the `launch` layer
* Contributes environment variables to configure the invoker if any configuration variables are defined. (Overrides anything from `func.yaml`)

## Configuration

| Environment Variable | Description |
|----------------------|-------------|
| `$BP_FUNCTION` | Configure the function to load. If the function lives in the default package: `<class>`. If the function lives in their own package: `<package>.<class>`. Defaults to `functions.Handler` |

## Getting started
To get started you'll need to create a directory where your function will be defined.

From within this directory we require a few files to properly detect this as a Java function:
* `func.yaml` (optional): We use this to configure the runtime environment variables.
  This buildpack makes use of `envs` and `options`. The keys `name` and `runtime` are required to maintain compatibility with Knative func cli, but are not used by this buildpack.
  See [Knative's func.yaml documentation](https://github.com/knative/func/blob/main/docs/reference/func_yaml.md)
  for more `func.yaml` information.
* `pom.xml` or `build.gradle`: These are used by the other Java buildpacks to compile your function.
* Java package in folder `src/main/java/functions`: This is the default location your function will be detected. If you do choose to use another package to store your functions, you will need to define where your function is located with the `BP_FUNCTION` configuration for the buildpack.

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
pack build <your_image_name_and_tag> --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:<version> --env BP_JVM_VERSION=17
```

Publish it to your registry:
```
docker push <your_image_name_and_tag>
```

## Liveness / Readiness Endpoints

The Java invoker contains a readiness/liveness endpoint that can be hit at `localhost:8080/actuator/health` by default. For more information, please read about the Spring Boot Actuator's [Kubernetes Probes](https://docs.spring.io/spring-boot/docs/2.3.0.RELEASE/reference/html/production-ready-features.html#production-ready-kubernetes-probes).

## Templates
If you want to quickly start writing your functions, take a look at the `templates/java` folder at the root of this repo.

## Examples
For some inspiration, take a look at the `samples/java` folder at the root of this repo.
