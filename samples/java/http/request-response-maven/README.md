# Java Sample - request-response

This is a sample function for how to create and send HTTP requests via our Java HTTP function.

## Prerequisites
    - Docker
    - Pack CLI

# Building - MAYBE DELETE

To build this sample, use the latest buildpack and build with `--env BP_FUNCTION=com.vmware.functions.Handler` set on the `pack` cli command.
## Usage
1. We want to first build the image:
    ```
    pack build request-response --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.1.0
    ```

1. After the image is successfully built we can run it in docker.
    ```
    docker run -p 8080:8080 --rm request-response
    ```

1. In a separate terminal we can send some POST requests to the function!
    ```
    REPLACE ME
    ```

1. You should see an expected result such as:
    ```
    REPLACE ME
    ```