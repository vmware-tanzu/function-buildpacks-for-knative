
# Function Buildpacks for Knative 

⚡ Build and deploy your functions fast -- a FaaS experience for HTTP and CloudEvents leveraging buildpacks.

Function Buildpacks for Knative (aka "Functions") brings functions as a programming model, to allow you to quickly build and deploy independent units of logic. Easily handle tasks such as asynchronous event reactions, cloud provider automations, and more. Soon, [func](https://github.com/knative-sandbox/kn-plugin-func) support will allow users to deploy Functions via CLI in a matter of seconds.

## Prerequisites

- [curl](https://curl.se/download.html) >= `7.79.0`
- [pack](https://buildpacks.io/docs/tools/pack/) >= `0.23.0`
- [func](https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/installing_cli.md) >= `0.21.2`
- [docker](https://docs.docker.com/get-docker/) (optional)

## Support Table

[![Test Builder](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-builder.yml/badge.svg)](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-builder.yml)
[![Test Buildpacks](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-buildpacks.yml/badge.svg)](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-buildpacks.yml)
[![Test Invokers](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-invokers.yml/badge.svg)](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-invokers.yml)
[![Test Templates](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-templates.yml/badge.svg)](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/test-templates.yml)

[![Dependency Review](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/dependency-review-action.yml/badge.svg)](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/dependency-review-action.yml)
[![Check Copyright and License](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/check-copyright-and-license.yml/badge.svg)](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/check-copyright-and-license.yml)
[![CodeQL](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/codeql-analysis.yml)

:warning: Currently in Alpha
| Language    | HTTP        | CloudEvents  |
| ----------- | ----------- | ------------ |
| Python      | Supported   | Supported    |
| Java        | Supported   | Supported    |
| NodeJS      | In-Progress | In-Progress  |
| .NET        | Planned     | Planned      |

## Getting Started

### Building

The buildpacks in this repository have been built, published, and included in the builder. 

The builder is also built and published as an image to `ghcr.io` -- to build an image from source, simply use the builder as shown below.

For example:
```
pack build <image_name> --path /path/to/function --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.1.0
```

You can get started by working with any of our scaffolded code found in `samples` or `templates` in the root directory.

To learn about local deployment and testing via `curl`, check out [DEPLOYING](DEPLOYING.md).

## Links

### Documentation

- [Java Function Buildpack](/buildpacks/java/README.md)
- [Python Function Buildpack](/buildpacks/python/README.md)

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
