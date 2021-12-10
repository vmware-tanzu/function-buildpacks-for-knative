
# Function Buildpacks for Knative

Enables a FaaS experience for Knative / Cloud Native Runtimes.

Will soon extend [func](https://github.com/knative-sandbox/kn-plugin-func) to create deployable functions via CLI.

## Pre-Requisites
- [Buildpack CLI](https://buildpacks.io/docs/tools/pack/)
- [Docker](https://docs.docker.com/get-docker/)
- [Knative](https://knative.dev/docs/) or [Cloud Native Runtimes](https://docs.vmware.com/en/Cloud-Native-Runtimes-for-VMware-Tanzu/1.0/tanzu-cloud-native-runtimes-1-0/GUID-cnr-overview.html)

## Currently Supported
* Python
* Java (Spring)

## Future
* NodeJS
* .NET
* Rust

## Usage
The buildpacks in this repository have been built, published, and included in the builder. 

The builder is also built and published as an image to `gcr.io` -- to build an image from source, simply use the builder as shown below.

For example:
```
pack build <image_name> --path /path/to/function --builder us.gcr.io/daisy-284300/kn-fn/builder:0.0.5
```

For more details you can refer to language-specific documentation:
* Python
    * [Buildpack/Function details](./buildpacks/python/README.md)
    * [Samples](./samples/python)
    * [Templates](./templates/python)
* Java (Alpha)
    * [Buildpack/Function details](./buildpacks/java/README.md)
    * [Samples](./samples/java)
    * [Templates](./templates/java)

To get started building the project, you can just run `make` in the root directory.
The Makefile includes calls to all the subdirectories. If you wish to run only a specific section, navigate to the subdirectory and run `make` or a related subcommand, e.g. `make buildpacks`, `make smoke-tests`.

## Documentation
Each subdirectory has a relevant README.md describing how to use its respective files.

## Contributing

The function-buildpacks-for-knative project team welcomes contributions from the community. Before you start working with function-buildpacks-for-knative, please
read our [Developer Certificate of Origin](https://cla.vmware.com/dco). All contributions to this repository must be
signed as described on that page. Your signature certifies that you wrote the patch or have the right to pass it on
as an open-source patch. For more detailed information, refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
* [BSD-2 License](LICENSE)

## Reporting Bugs or Vulnerabilities
* [Bugs, Issues, Missing Features](https://github.com/vmware-tanzu/function-buildpacks-for-knative/issues/)
* [Only Vulnerabilities](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/SECURITY.md)
