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
)

func TestHelloWorldHTTP(t *testing.T) {
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

func TestHelloWorldCloudEvents(t *testing.T) {
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

			path:             "/hello",
			expectedResponse: "Hello World!",
		},
		{
			name: "Python",
			tag:  "python",

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

			client, err := cloudevents.NewClientHTTP()
			if err != nil {
				t.Error(err)
			}

			event := cloudevents.NewEvent()
			event.SetSource("url")
			event.SetType("example.type")
			event.SetData(cloudevents.ApplicationJSON, "Hello World!")

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
