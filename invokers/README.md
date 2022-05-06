# Invokers

These language-specific invokers are built and bundled into the buildpack (found at the root `buildpacks/` folder) via the `buildpack.toml` file found under each buildpacks' respective language. These invokers are used to invoke the function in the language of your choice.

## Development Workflow

> :warning: **Change your project's registry to a development one**
>
> Navigate to `rules.mk` in the **root directory** and replace the registry with your own. By default it is set to `us.gcr.io/daisy-284300`, which is for VMware developers only, and is not a true development environment. Ensure you have sufficient privileges and re-authenticate with your CLI of choice, such as using `gcloud auth login`.

### Python Local Development

To test your Python invoker changes locally, you can create a virtual environment via `python3 -m venv venv` and running `pip install . ` in the Python invoker's directory. Then, navigate to a template or test folder of your choice, and run `python3 -m pyfunc start` to launch the Python invoker locally upon your function's code.

### Testing Invoker Changes

If you are a developer making changes to these invokers, there are several generic steps you must follow:

After making your changes to your Invoker's files:

1) Use the provided `Makefile`'s `publish` command to push your changes to a GCR repo (e.g. `make invoker.<language>.publish`)
2) Note down that link and the `SHA` provided in the newly generated `out` folder
3) Then, modify the respective language's `buildpacks/<language>/buildpack.toml`, and then cut a new release of the buildpack to GCR
4) Update `builder/builder.toml` to finally push all these changes to the builder
5) Locally build your function with this new builder

All of these processes can be explored with a further dive into the `Makefile`'s of each component.
