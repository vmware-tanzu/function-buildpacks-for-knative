# Java Sample - request-response

This is a sample function for how to create and send HTTP requests via our Java HTTP function.

## Prerequisites
    - Docker
    - Pack CLI

## Usage
1. We want to first build the image:
    ```
    pack build request-response --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.1.0 --env BP_FUNCTION=functions.Handler
    ```

1. After the image is successfully built we can run it in docker.
    ```
    docker run -p 8080:8080 --rm request-response
    ```

1. In a separate terminal we can send some POST requests to the function!
    ```
    # Sample request to convert 100 degrees celsius to fahrenheit. 
    curl -X POST localhost:8080 -H "Content-Type: application/json" -d '{"celsius":100}'
    ```
   The result should be 212.0 
