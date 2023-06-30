To create the container with the Pack CLI:
```
pack build s3lambda --builder paketobuildpacks/builder:0.3.50-base --post-buildpack ghcr.io/vmware-tanzu/function-buildpacks-for-knative/python-buildpack:1.1.2 --env BP_FUNCTION=func.main
```
