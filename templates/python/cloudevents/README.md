# Python CloudEvents Template

## Getting Started

## Testing

If you'd like to test this template, you may use this CloudEvent saved as a `.json` file:

```
{
    "specversion" : "1.0",
    "type" : "com.github.pull_request.opened",
    "source" : "https://github.com/cloudevents/spec/pull",
    "subject" : "123",
    "id" : "A234-1234-1234",
    "time" : "2018-04-05T17:31:00Z",
    "comexampleextension1" : "value",
    "comexampleothervalue" : 5,
    "datacontenttype" : "text/plain",
    "data" : "helloworld"
}
```

After [deploying](./DEPLOYING.md) your function as an image, you can test with:

`curl -X POST -H "Content-Type: application/cloudevents+json" -d @ce2.json http://localhost:8080`
