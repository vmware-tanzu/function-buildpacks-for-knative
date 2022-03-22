package tests

import (
	"bytes"
	"context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"io/ioutil"

	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestSmokeHTTP(t *testing.T) {
	baseImage := "kn-fn-test/helloworld"
	cases := []struct {
		name string
		tag  string

		methodType       string
		contentType      string
		path             string
		expectedResponse string
	}{
		{
			name: "Java",
			tag:  "java",

			methodType:       http.MethodPost,
			path:             "/hello",
			expectedResponse: "Hello World!",
		},
		{
			name: "Python GET",
			tag:  "python",

			methodType:       http.MethodGet,
			path:             "/",
			expectedResponse: "Hello World!",
		},
		{
			name: "Python POST",
			tag:  "python",

			methodType:       http.MethodPost,
			path:             "/",
			expectedResponse: "Hello World!",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			image := fmt.Sprintf("%s:%s", baseImage, c.tag)
			cleanup, err := runTestContainer(image)
			if err != nil {
				t.Error(err)
				return
			}
			defer cleanup()

			url := fmt.Sprintf("http://127.0.0.1:8080/%s", strings.TrimLeft(c.path, "/"))

			var resp *http.Response
			switch c.methodType {
			case http.MethodGet:
				resp, err = http.Get(url)
				if err != nil {
					t.Error(err)
					return
				}
			case http.MethodPost:
				ct := c.contentType
				if ct == "" {
					ct = "application/json"
				}
				resp, err = http.Post(url, ct, bytes.NewBufferString(""))
				if err != nil {
					t.Error(err)
					return
				}
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
				return
			}

			bs := string(body)
			if bs != c.expectedResponse {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, bs)
			}
		})
	}
}

func TestSmokeCloudEvents(t *testing.T) {
	baseImage := "kn-fn-test/echo-ce"
	cases := []struct {
		name string
		tag  string

		path             string
		data             string
		expectedResponse string
	}{
		{
			name: "Java",
			tag:  "java",

			path:             "/",
			data:             "java test data",
			expectedResponse: "java test data",
		},
		{
			name: "Python",
			tag:  "python",

			path:             "/",
			data:             "python test data",
			expectedResponse: "python test data",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			image := fmt.Sprintf("%s:%s", baseImage, c.tag)
			cleanup, err := runTestContainer(image)
			if err != nil {
				t.Error(err)
				return
			}
			defer cleanup()

			url := fmt.Sprintf("http://127.0.0.1:8080/%s", strings.TrimLeft(c.path, "/"))

			client, err := cloudevents.NewClientHTTP()
			if err != nil {
				t.Error(err)
			}

			event := cloudevents.NewEvent()
			event.SetSource("url")
			event.SetID(uuid.New().String())
			event.SetType("example.type")
			event.SetData(cloudevents.TextPlain, c.data)

			ctx := cloudevents.ContextWithTarget(context.Background(), url)
			reqEvent, result := client.Request(ctx, event)

			if cloudevents.IsUndelivered(result) {
				t.Error(err)
			}

			// Extra check due to odd behavior in CloudEvents Go SDK: github.com/cloudevents/sdk-go/blob/1170e89edb9b504a806f2c6a26563c3c26b68276/v2/client/client.go#L178
			if cloudevents.IsNACK(result) {
				t.Error(err)
				t.Skip()
			}

			reqEventData := string(reqEvent.Data())
			if reqEventData != c.expectedResponse {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, reqEventData)
			}
		})
	}
}

func TestTemplatesHTTP(t *testing.T) {
	baseImage := "kn-fn-test/template-http"
	jsonData := []byte(`{"firstName":"John", "lastName":"Doe"}`)
	cases := []struct {
		name string
		tag  string

		methodType       string
		contentType      string
		path             string
		expectedResponse string
	}{
		{
			name: "Java HTTP Gradle",
			tag:  "java-http-gradle",

			methodType:       http.MethodPost,
			data:             jsonData,
			path:             "/hire",
			expectedResponse: "Hello World!",
		},
		{
			name: "Java HTTP Maven",
			tag:  "java-http-maven",

			methodType:       http.MethodPost,
			data:             jsonData,
			path:             "/hire",
			expectedResponse: "Hello World!",
		},
		{
			name: "Python GET",
			tag:  "python-http",

			methodType:       http.MethodGet,
			path:             "/",
			expectedResponse: "Hello World!",
		},
		{
			name: "Python POST",
			tag:  "python-http",

			methodType:       http.MethodPost,
			path:             "/",
			expectedResponse: "Hello World!",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			image := fmt.Sprintf("%s:%s", baseImage, c.tag)
			cleanup, err := runTestContainer(image)
			if err != nil {
				t.Error(err)
				return
			}
			defer cleanup()

			url := fmt.Sprintf("http://127.0.0.1:8080/%s", strings.TrimLeft(c.path, "/"))

			var resp *http.Response
			switch c.methodType {
			case http.MethodGet:
				resp, err = http.Get(url)
				if err != nil {
					t.Error(err)
					return
				}
			case http.MethodPost:
				ct := c.contentType
				if ct == "" {
					ct = "application/json"
				}

				if c.data != nil {
					resp, err = http.Post(url, ct, bytes.NewBuffer(jsonData))

				} else {
					resp, err = http.Post(url, ct, bytes.NewBufferString(""))
				}

				if err != nil {
					t.Error(err)
					return
				}
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
				return
			}

			bs := string(body)
			if bs != c.expectedResponse {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, bs)
			}
		})
	}
}

func TestTemplatesCloudEvents(t *testing.T) {
	baseImage := "kn-fn-test/template-ce"
// 	jsonData := byte[](`{
//     "specversion" : "1.0",
//     "type" : "org.springframework",
//     "source" : "https://spring.io/",
//     "id" : "A234-1234-1234",
//     "datacontenttype" : "application/json",
//     "data": {
//         "firstName": "John",
//         "lastName": "Doe"
//     }
// }`)
	cases := []struct {
		name string
		tag  string

		path             string
		data             string
		expectedResponse string
	}{
		{
			name: "Java CloudEvents Gradle",
			tag:  "java-cloudevents-gradle",

			path:             "/hire",
			data:             "REPLACEME",
			expectedResponse: "java test data",
		},
		{
			name: "Java CloudEvents Maven",
			tag:  "java-cloudevents-maven",

			path:             "/hire",
			data:             "REPLACEME",
			expectedResponse: "java test data",
		},
		{
			name: "Python CloudEvents",
			tag:  "python-cloudevents",

			path:             "/",
			data:             "python test data",
			expectedResponse: "python test data",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			image := fmt.Sprintf("%s:%s", baseImage, c.tag)
			cleanup, err := runTestContainer(image)
			if err != nil {
				t.Error(err)
				return
			}
			defer cleanup()

			url := fmt.Sprintf("http://127.0.0.1:8080/%s", strings.TrimLeft(c.path, "/"))

			client, err := cloudevents.NewClientHTTP()
			if err != nil {
				t.Error(err)
			}

			event := cloudevents.NewEvent()
			event.SetSource("url")
			event.SetID(uuid.New().String())
			event.SetType("example.type")
			event.SetData(cloudevents.TextPlain, c.data)

			ctx := cloudevents.ContextWithTarget(context.Background(), url)
			reqEvent, result := client.Request(ctx, event)

			if cloudevents.IsUndelivered(result) {
				t.Error(err)
			}

			// Extra check due to odd behavior in CloudEvents Go SDK: github.com/cloudevents/sdk-go/blob/1170e89edb9b504a806f2c6a26563c3c26b68276/v2/client/client.go#L178
			if cloudevents.IsNACK(result) {
				t.Error(err)
				t.Skip()
			}

			reqEventData := string(reqEvent.Data())
			if reqEventData != c.expectedResponse {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, reqEventData)
			}
		})
	}
}
