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
2) `make builder.image.local`

then you may try the new `pack build` command outputted from that result. Include `--pull-policy if-not-present`.  This will use local images first and pull any missing dependencies. 

Please note that these `uri`s should remain unchanged except for a version number bump for RELEASE builds, so your PR should contain the necessary changes to `buildpacks` and `invokers`, but not necessarily to this directory.

### Cut a Builder Release

1. If you modified the invokers or buildpacks, following the [Cut a Buildpack Release](../buildpacks/README.md) instructions.
2. If you haven't already, disable Github Branch Protection rules as described in the Buildpack Release Instructions. 
3. Choose [Create Builder Release](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/create-builder-release.yaml)
4. Click "Run Workflow" and choose the type of release.  This will create a commit that updates builder/VERSION.
5. Search through the project files for `pack build` commands that reference the version of the builder. Update the builder version to the new builder's version.
6. Re-enable Branch Protection rules. 
