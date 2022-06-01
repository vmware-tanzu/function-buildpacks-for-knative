# Builder -- function-buildpacks-for-knative

The builder aggregates the sum of our layers, by maintaining a collection of `uri`'s that contain the buildpack images and more

In abstract, the concept of layers may be viewed as (from lowest level to highest): invoker -> buildpack -> builder.

## Testing your changes

Once you have obtained your buildpack's `uri` that you wish to update (see: [buildpacks/README.md](https://github.com/vmware-tanzu/function-buildpacks-for-knative/tree/main/buildpacks/README.md) for more information), open the `builder.toml` file in this directory to make your changes.

You will find under the `buildpacks` variable a map of { `id` => `uri` }. Please update this with your development buildpack layer for testing purposes.

To test your newly edited builder, run `make builder.clean`, then `make builder`), and then you may try the new `pack build` command outputted from that result.

Please note that these `uri`s should remain unchanged except for a version number bump for RELEASE builds, so your PR should contain the necessary changes to `buildpacks` and `invokers`, but not necessarily to this directory.

