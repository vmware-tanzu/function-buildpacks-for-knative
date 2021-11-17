package tests

import (
	"bytes"
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"io/ioutil"
	"log"

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

			methodType:       "example.type",
			path:             "/hello",
			expectedResponse: "Hello World!",
		},
		{
			name: "Python",
			tag:  "python",

			methodType:       "example.type",
			path:             "/",
			expectedResponse: "Hello World!",
		},
	}

	// Create a CloudEvents client
	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Somewhere in here, a timeout occurs, preventing tests from passing
			// Double check this logic
			image := fmt.Sprintf("%s:%s", baseImage, c.tag)
			cleanup, err := runTestContainer(image) // Container probably needs to be modified to launch a CE http client
			if err != nil {
				t.Error(err)
				return
			}
			defer cleanup()

			url := fmt.Sprintf("http://127.0.0.1:8080/%s", strings.TrimLeft(c.path, "/")) // unsure if this is still valid

			// Send a CloudEvent
			event := cloudevents.NewEvent()
			event.SetSource("url")
			event.SetType("example.type")
			event.SetData(cloudevents.ApplicationJSON, "Hello World!")

			// Set a target.
			ctx := cloudevents.ContextWithTarget(context.Background(), url)

			// Send that Event.
			if result := client.Send(ctx, event); cloudevents.IsUndelivered(result) {
				log.Fatalf("failed to send, %v", result)
			}

			// Receive a CloudEvent
			log.Fatal(client.StartReceiver(context.Background(), func(e cloudevents.Event) {
				bs := string(e.Data()[:])
				if bs != c.expectedResponse {
					t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, bs)
				}
			}))

		})
	}
}
