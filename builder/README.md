# Builder -- function-buildpacks-for-knative

The builder aggregates the sum of our layers, by maintaining a collection of `uri`'s that contain the buildpack images and more

In abstract, the concept of layers may be viewed as (from lowest level to highest): invoker -> buildpack -> builder.

## Testing your changes

> :warning: **Change your project's registry to a development one**
>
> Navigate to `rules.mk` in the **root directory** and replace the registry with your own. By default it is set to `us.gcr.io/daisy-284300`, which is for VMware developers only, and is not a true development environment. Ensure you have sufficient privileges and re-authenticate with your CLI of choice, such as using `gcloud auth login`.

Once you have obtained your buildpack's `uri` that you wish to update (see: [buildpacks/README.md](https://github.com/vmware-tanzu/function-buildpacks-for-knative/tree/main/buildpacks/README.md) for more information), open the `builder.toml` file in this directory to make your changes.

You will find under the `buildpacks` variable a map of { `id` => `uri` }. Please update this with your development buildpack layer for testing purposes.

To test your newly edited builder, run:
1) `make builder.clean`
2) `make builder`
3) `make builder.publish`

then you may try the new `pack build` command outputted from that result.

Please note that these `uri`s should remain unchanged except for a version number bump for RELEASE builds, so your PR should contain the necessary changes to `buildpacks` and `invokers`, but not necessarily to this directory.

