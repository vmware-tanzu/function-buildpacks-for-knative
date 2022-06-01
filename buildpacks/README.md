# Buildpacks -- function-buildpacks-for-knative

You may find your language-specific buildpacks here and edit them.

## Testing your changes

To test your changes, ensure you have taken the necessary steps to access your repository with read/write access. This includes updating your repository type and location, found at [rules.mk](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/rules.mk).

After, you can run `make <language>-buildpack.publish` to publish your language's new buildpack layer. (e.g. `make python-buildpack.publish`)

With the new registry's URL obtained by the `make` command's output, you can now update [builder/builder.toml](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/builder/builder.toml) to use your development buildpack.

Please see the builder's [README.md](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/builder/README.md) for more information.
