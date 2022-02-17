# Java CloudEvents Template

## Getting Started

Navigate to `Hire.java` and replace the code in the function body as desired.

## Testing

If you'd like to test this template, you may use this CloudEvent saved as `cloudevent.json`:

```
{
    "specversion" : "1.0",
    "type" : "org.springframework",
    "source" : "https://spring.io/",
    "id" : "A234-1234-1234",
    "datacontenttype" : "application/json",
    "data": {
        "firstName": "John",
        "lastName": "Doe"
    }
}
```

After [deploying](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/DEPLOYING.md) your function as an image, you can test with:

```
curl -X POST -H "Content-Type: application/cloudevents+json" -d @cloudevent.json http://localhost:8080
```
