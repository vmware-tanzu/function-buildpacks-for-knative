To deploy this, first create the container with the buildpack cli
```
pack build <your_image_name_and_tag> --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.1.0
```

Publish it to your registry:
```
docker push <your_image_name_and_tag>
```