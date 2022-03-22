//go:build template
// +build template

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

func TestPythonHTTP(t *testing.T) {
	baseImage := "kn-fn-test/template-http"
	cases := []struct {
		name string
		tag  string

		methodType       string
		contentType      string
		path             string
		expectedResponse string
	}{
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

func TestJavaHTTP(t *testing.T) {
	baseImage := "kn-fn-test/template-http"
	jsonData := []byte(`{"firstName":"John", "lastName":"Doe"}`)
	cases := []struct {
		name string
		tag  string

		methodType       string
		contentType      string
		data             []byte
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

				resp, err = http.Post(url, ct, bytes.NewBuffer(jsonData))

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

func TestPythonCloudEvents(t *testing.T) {
	baseImage := "kn-fn-test/template-ce"
	jsonData := []byte(`{
		"specversion" : "1.0",
		"type" : "org.springframework",
		"source" : "https://spring.io/",
		"id" : "A234-1234-1234",
		"datacontenttype" : "application/json",
		"data": {
			"firstName": "John",
			"lastName": "Doe"
		}
	}`)
	cases := []struct {
		name string
		tag  string

		path             string
		data             []byte
		expectedResponse string
	}{
		{
			name: "Python CloudEvents",
			tag:  "python-cloudevents",

			path:             "/",
			data:             jsonData,
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

func TestJavaCloudEvents(t *testing.T) {
	baseImage := "kn-fn-test/template-ce"
	jsonData := []byte(`{
	    "specversion" : "1.0",
	    "type" : "org.springframework",
	    "source" : "https://spring.io/",
	    "id" : "A234-1234-1234",
	    "datacontenttype" : "application/json",
	    "data": {
	        "firstName": "John",
	        "lastName": "Doe"
	    }
	}`)
	cases := []struct {
		name string
		tag  string

		path             string
		data             []byte
		expectedResponse string
	}{
		{
			name: "Java CloudEvents Gradle",
			tag:  "java-cloudevents-gradle",

			path:             "/hire",
			data:             jsonData,
			expectedResponse: "java test data",
		},
		{
			name: "Java CloudEvents Maven",
			tag:  "java-cloudevents-maven",

			path:             "/hire",
			data:             jsonData,
			expectedResponse: "java test data",
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
