# Invokers

* This assumes that you have permissions to access VMware's GCR -- if you do not, please manually build these steps or point to your own cloud repository.
* Don't forget to `gcloud auth login` if you get a 401

## Development Workflow

After making your changes to your Invoker's files:

1) Use the provided `Makefile`'s `publish` command to push your changes to a GCR repo (e.g. `make invoker.python.publish`)
2) Note down that link and the `SHA` provided in the newly generated `out` folder
3) Then, modify the respective language's `buildpacks/<language>/buildpack.toml`, and then cut a new release of the buildpack to GCR
4) Update `builder/builder.toml` to finally push all these changes to the builder
5) Locally build your function with this new builder