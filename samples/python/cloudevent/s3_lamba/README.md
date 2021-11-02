To deploy this, first create the container with the buildpack cli
```
pack build <your_image_name_and_tag> --builder us.gcr.io/daisy-284300/kn-fn/builder:0.0.1
```

Publish it to your registry:
```
docker push <your_image_name_and_tag>
```