# Java Sample - request-response

This is a sample function for how to create and send HTTP requests via our Java HTTP function.

## Prerequisites

- Docker
- Pack CLI

## Usage

1.  Build the image:

    ```
    pack build request-response --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.1.0 --env BP_FUNCTION=com.vmware.functions.Handler
    ```

1. Run it in Docker:

    ```
    docker run -p 8080:8080 --rm request-response
    ```

1. In a separate terminal, send some requests to the function!

    1. Issue a GET to the `uppercase()` function: 
        ```
        curl localhost:8080/uppercase/lowercase
        ```
        Expected result: `LOWERCASE`

    1. Issue a GET to the `convert()` function:
        ```
        curl http://localhost:8080/convert/%7B%22celsius%22:100%7D
        ```
        Note that the URI is `/convert/{"celsius":100}`. For the curl some characters are encoded. 
        Expected result: `212.0`

    1. Issue a POST to the `convert()` function:
        ```
        curl -X POST localhost:8080/convert -H "Content-Type: application/json" -d '{"celsius":100}'
        ```
        Expected result: `212.0`

## Health Endpoints

The Spring actuator is available.  Check the health endpoint: 
```
curl localhost:8080/actuator/health
```
Expected result: `{"status":"UP"}`
