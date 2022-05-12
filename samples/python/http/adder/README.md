# Adder

This example will take two number and add them up only if the requestor is authenticated.

# Running the sample
1. We want to first build the image:
    ```
    pack build adder --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.0.10
    ```
    * `VERSION` is the version of the builder to use.

1. After the image is successfully built we can run it in docker.
    ```
    docker run -p 8080:8080 --rm adder
    ```

1. In a separate terminal we can send some requests to the function!
    ```
    # Sample request to do 65+6 as user "admin"
    curl localhost:8080 -F username=admin -F password=supersecure -F first=65 -F second=6
    ```
    Accepted form key and values:
    * `username`
        * Valid users: `admin`, `asu`, `someone`
        * If this value is not specified or incorrect, it will fail with HTTP status 401
    * `password`
        * Passwords for the above users in the order: `supersecure`, `mypassword`, `123qwe`
        * If this value is not specified or incorrect, it will fail with HTTP status 401
    * `first`
        * If not specified, default is `0`
        * Anything other than integers will cause an HTTP status 500
    * `second`
        * If not specified, default is `0`
        * Anything other than integers will cause an HTTP status 500
