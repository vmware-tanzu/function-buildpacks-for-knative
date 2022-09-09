# Buildpacks -- function-buildpacks-for-knative

This folder contains language-specific buildpacks. 

## Running locally

### Prerequisite
Before you can build a local buildpack, you'll need the invoker files created.  
Run `make invokers.<language>`

### Building
To make a local buildpack image, run `make buildpacks.<language>.images.local`.

### Running
To run your newly built local buildpack, use the `pack` CLI. 

_Java example:_ 
```
cd <path-to-java-function>
pack build \
  -b gcr.io/paketo-buildpacks/java:7.0.0 \
  -b <local-image> \
  --builder paketobuildpacks/builder-jammy-buildpackless-base:latest \
  --verbose --clear-cache --pull-policy if-not-present \
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