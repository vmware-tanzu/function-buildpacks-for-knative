
# Function Buildpacks for Knative 

Develop with a FaaS experience for HTTP and CloudEvents leveraging Cloud Native Buildpacks.

Will soon extend [func](https://github.com/knative-sandbox/kn-plugin-func) to create deployable functions via CLI.

## Prerequisites
- [curl](https://curl.se/download.html) >= `7.79.0`
- [pack](https://buildpacks.io/docs/tools/pack/) >= `0.23.0`
- [func](https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/installing_cli.md) >= `0.21.2`
- [docker](https://docs.docker.com/get-docker/) (optional)

## Support Table
:warning: Currently in Alpha
| Language    | HTTP        | CloudEvents  |
| ----------- | ----------- | ------------ |
| Python      | Supported   | Supported    |
| Java        | Supported   | Supported    |
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

To learn about local deployment and testing via `curl`, check out [DEPLOYING](DEPLOYING.md).

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
