
# Function Buildpacks for Knative

Develop with a FaaS experience for HTTP and CloudEvents leveraging Cloud Native Buildpacks.

Will soon extend [func](https://github.com/knative-sandbox/kn-plugin-func) to create deployable functions via CLI.

## Pre-Requisites
- [curl](https://curl.se/download.html)
- [pack](https://buildpacks.io/docs/tools/pack/) >= `0.23.0`
- [func](https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/installing_cli.md)
- [docker](https://docs.docker.com/get-docker/)

## Support Table
| Language    | HTTP        | CloudEvents  |
| ----------- | ----------- | ------------ |
| Python      | Supported   | Experimental |
| Java        | Supported   | Experimental |
| NodeJS      | Planned     | Planned      |
| .NET        | Planned     | Planned      |

## Getting Started

### Building
The buildpacks in this repository have been built, published, and included in the builder. 

The builder is also built and published as an image to `ghcr.io` -- to build an image from source, simply use the builder as shown below.

For example:
```
pack build <image_name> --path /path/to/function --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.0.8
```

You can get started by working with any of our scaffolded code found in `samples` or `templates` in the root directory.


### Deploying

This assumes you have Docker Desktop properly installed and running.

With Docker Desktop running, authenticated, and the ports (default `8080`) available:

```
docker run -it --rm -p 8080:8080 sample-java
```

### Testing
After deploying your function, you can interact with our templates by doing:
- Single function definition: `curl -X POST localhost:8080`
- Multiple function definitions: `curl -H "Content-Type: application/json" -X POST localhost:8080/hello`
- - where `hello` as a path invokes your function's definition

For more details you can refer to language-specific documentation:
* Python
    * [Buildpack](./buildpacks/python/README.md)
    * [Samples](./samples/python)
    * [Templates](./templates/python)
* Java / Spring
    * [Buildpack](./buildpacks/java/README.md)
    * [Samples](./samples/java)
    * [Templates](./templates/java)

To get started building the project, you can just run `make` in the root directory.
The Makefile includes calls to all the subdirectories. If you wish to run only a specific section, navigate to the subdirectory and run `make` or a related subcommand, e.g. `make buildpacks`, `make smoke-tests`.

## Links

### Contributing

The function-buildpacks-for-knative project team welcomes contributions from the community. Before you start working with function-buildpacks-for-knative, please
read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be
signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on
as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

### License
* [BSD-2 License](LICENSE)

### Reporting Bugs or Vulnerabilities
* [Bugs, Issues, Missing Features](https://github.com/vmware-tanzu/function-buildpacks-for-knative/issues/)
* [Only Vulnerabilities](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/SECURITY.md)
