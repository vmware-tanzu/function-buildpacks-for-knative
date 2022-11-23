# Invokers -- function-buildpacks-for-knative

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

1) Use the provided `Makefile`'s `publish` command to push your changes to a GCR repo (e.g. `make invokers.<language>.publish`)
2) Note down that link and the `SHA` provided in the newly generated `out` folder
3) Then, modify the respective language's `buildpacks/<language>/buildpack.toml`, and then cut a new release of the buildpack to GCR (`make buildpacks.<language>.publish`)
4) Update `builder/builder.toml` to finally push all these changes to the builder
5) Locally build your function with this new builder (`make builder.clean` if needed, then `make builder`)

All of these processes can be explored with a further dive into the `Makefile`'s of each component.

### Cut a Buildpack Release After Modifying Invokers

To build new versioned invokers:  

1. Test (as discussed above) and commit your changes to the main branch. Observe successful checks on the PR. 
1. Via the Github web UI navigate to Actions and  
   choose [Create Invoker Release](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/create-invoker-release.yaml)
1. Click "Run Workflow", choose the language and the type of release. This will result in the action making a commit that changes the invoker's VERSION file. 

To build a new versioned buildpack: 
1. If necessary, update the cpe versions in buildpacks/<language>/ytt/dependency-metadata.yaml where <language> is java or python. 
1. Choose [Create Buildpack Release](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/create-buildpack-release.yaml)
1. Click "Run Workflow", choose the language and type of release. 
   This will create a commit that updates buildpacks/<language>/VERSION and buildpacks/<language>/buildpack.toml

Update the builder: 
1. Choose [Create Builder Release](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/create-builder-release.yaml)
1. Click "Run Workflow" and choose the type of release.  This will create a commit that updates builder/VERSION. 
1. Search through documentation files for `pack build` commands that reference the version of the builder. Update
   the builder version to the new builder's version. 