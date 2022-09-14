# Deploying

## Building your function

You can build your function using our provided builder, which already includes buildpacks and an invoker layer:
```
pack build my-function --path . --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.1.0
```
Where `my-function` is the name of your runnable function image, later used by Docker.

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
