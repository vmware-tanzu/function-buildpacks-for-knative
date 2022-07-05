# Buildpacks -- function-buildpacks-for-knative

You may find your language-specific buildpacks here and edit them.

## Testing your changes

> :warning: **Change your project's registry to a development one**
>
> Navigate to `rules.mk` in the **root directory** and replace the registry with your own. By default it is set to `us.gcr.io/daisy-284300`, which is for VMware developers only, and is not a true development environment. Ensure you have sufficient privileges and re-authenticate with your CLI of choice, such as using `gcloud auth login`.

To test your changes, ensure you have taken the necessary steps to access your repository with read/write access. This includes updating your repository type and location, found at [rules.mk](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/rules.mk).

After, you can run `make buildpacks.<language>.images.publish` to publish your language's new buildpack. (e.g. `make buildpacks.python.images.publish`)

With the new registry's URL obtained by the `make` command's output, you can now update [builder/builder.toml](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/builder/builder.toml) to use your development buildpack.

Please see the builder's [README.md](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/builder/README.md) for more information.
