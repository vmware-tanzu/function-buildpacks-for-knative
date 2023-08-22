To deploy this, first create the container with the buildpack cli
```
pack build <your_image_name_and_tag> --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.4.3 --env BP_FUNCTION=main.main
```

Publish it to your registry:
```
docker push <your_image_name_and_tag>
```