# Buildpacks -- function-buildpacks-for-knative

This folder contains language-specific buildpacks. 

## Testing your buildpacks

### Prerequisite
Before you can build a local buildpack, you'll need the invoker files created.  
Run `make invokers.<language>`

### Building
To make a buildpack locally run `make buildpacks.<language>.images.local`.  The output of this 
is a built image in your local registry. 

### Testing
To test your newly built local buildpack use the pack command. 

Java example: 
```
cd <path-to-java-function>
pack build \
  -b gcr.io/paketo-buildpacks/java:7.0.0 \
  -b <local-image> \
  --builder paketobuildpacks/builder-jammy-buildpackless-base:latest \
  --verbose --clear-cache --pull-policy if-not-present \
  <output-image>
```

Python example: 
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
