# Buildpacks -- function-buildpacks-for-knative

This folder contains language-specific buildpacks. 

## Running locally

### Building
To make a local buildpack image, run `make buildpacks.<language>.images.local`.

### Running
To run your newly built local buildpack, use the `pack` CLI. 

_Java example:_ 
```
cd <path-to-java-function>
pack build \
  -b gcr.io/paketo-buildpacks/java:8.3.0 \
  -b <local-image> \
  --builder paketobuildpacks/builder-jammy-buildpackless-base:latest \
  --verbose --clear-cache --pull-policy if-not-present \
  --env BP_JVM_VERSION=17 \
  <output-image>
```

_Python example:_ 
```
cd <path-to-python-function>
pack build \
  -b gcr.io/paketo-buildpacks/python:2.0.0 \
  -b <local-image> \
  --builder paketobuildpacks/builder-jammy-buildpackless-base:latest \
  --verbose --clear-cache --pull-policy if-not-present \
  <output-image>
```

where `<local-image>` is the output of `make buildpacks.<language>.images.local` 
and `<output-image>` is the function image you are building. 

## Testing

To run unit tests for a buildpack, run `make buildpacks.<language>.tests`. See the [testing documentation](/tests/README.md) to run other kinds of tests related to buildpacks.

### Cut a Buildpack Release

First, create releases of the invoker(s):

1. Test (as discussed above) and commit your changes to the main branch. Observe successful checks on the PR.
2. Enable temporary admin access to the github repo by going to [Upstream Contrib](https://upstreamcontrib.eng.vmware.com/oss/#/login) choosing the
   temporary admin access request in a drop-down.  Then go into the Github repo Settings, Branches, Edit the Branch Protection rule for the main branch
   and turn off both "Require a pull request before merging" and "Restrict who can push to matching branches".
3. Via the Github web UI navigate to Actions and choose [Create Invoker Release](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/create-invoker-release.yaml)
4. Click "Run Workflow", choose the language and the type of release. This will result in the action making a commit that changes the invoker's VERSION file.

Then, create the buildpack release(s): 

1. If necessary, update the Common Platform Enumerations (cpe) versions in buildpacks/python/ytt/dependency-metadata.yaml and/or buildpacks/java/ytt/dependency-metadata.yaml.
   Find the data necessary to specify cpe entries by downloading the dependency asset from the Github invoker release created above.
   Entries should follow the pattern from [NIST](https://nvd.nist.gov/products/cpe/search).
2. Choose [Create Buildpack Release](https://github.com/vmware-tanzu/function-buildpacks-for-knative/actions/workflows/create-buildpack-release.yaml)
3. Click "Run Workflow", choose the language and type of release.
   This will create a commit that updates `buildpacks/<language>/VERSION` and `buildpacks/<language>/buildpack.toml`. 
   Automation will pull these buildpacks into the Tanzu buildpacks.
4. Re-enable Branch Protection rules.
