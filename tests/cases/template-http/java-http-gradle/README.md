# Java HTTP Template

## Getting Started

Navigate to `Hire.java` and replace the code in the function body as desired.

## Testing

After [deploying](https://github.com/vmware-tanzu/function-buildpacks-for-knative/blob/main/DEPLOYING.md) your function as an image, you can test with:
```
curl -w'\n' localhost:8080/hire \
 -H "Content-Type: application/json" \
 -d '{"firstName":"John", "lastName":"Doe"}' -i
```

## Notes

Additionally you can make a CloudEvent request and receive the appropriate CloudEvent response with the following:
```
curl -w'\n' localhost:8080/hire \
 -H "ce-id: 0001" \
 -H "ce-specversion: 1.0" \
 -H "ce-type: hire" \
 -H "ce-source: spring.io/spring-event" \
 -H "Content-Type: application/json" \
 -d '{"firstName":"John", "lastName":"Doe"}' -i
```
