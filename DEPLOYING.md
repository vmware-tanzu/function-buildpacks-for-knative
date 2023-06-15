# Deploying

## Building your function

You can build your function using our provided builder, which already includes buildpacks and an invoker layer:
```
pack build my-function --path . --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.4.2  --env BP_FUNCTION=path.function --env BP_JVM_VERSION=17
```
Where:  
  * `my-function` is the name of your runnable function image, later used by Docker.
  * `path` is the name of the file or package where the function resides. 
  * `function` is the name of the method or function. 

Examples: 
  * Python: BP_FUNCTION=func.main. `func` is the name of the .py file. main is the `method`.
  * Java: BP_FUNCTION=function.Handler. `function` the package. `Handler` is the class that implements Function.

> Note: The invoker layer for Java functions utilizes Spring 3.0 which requires Java 17+. Add `--env BP_JVM_VERSION=17` to your `pack build` command to influence the buildpacks to add JRE version 17 to your resulting image instead of the default version 11. 

## Local Deployment

### Docker

This assumes you have Docker Desktop properly installed and running.

With Docker Desktop running, authenticated, and the ports (default `8080`) available:

```
docker run -it --rm -p 8080:8080 my-function
```

## Testing
After deploying your function, you can interact with our templates by doing:
- Single function definition: `curl -X POST localhost:8080`
- Multiple function definitions: `curl -H "Content-Type: application/json" -X POST localhost:8080/hello`
  - where `hello` as a path invokes your function's definition

With our templates, you should see some HTML or sample text returned indicating a success.
