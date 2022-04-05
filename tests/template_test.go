// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

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

			actualResponse := string(body)
			if actualResponse != c.expectedResponse {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, actualResponse)
			}
		})
	}
}

func TestJavaHTTP(t *testing.T) {
	baseImage := "kn-fn-test/template-http"
	jsonData := []byte(`{"firstName":"John","lastName":"Doe"}`)
	expectedData := `{"firstName":"John","lastName":"Doe"}`
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
			contentType:      "application/json",
			data:             jsonData,
			path:             "/hire",
			expectedResponse: expectedData,
		},
		{
			name: "Java HTTP Maven",
			tag:  "java-http-maven",

			methodType:       http.MethodPost,
			contentType:      "application/json",
			data:             jsonData,
			path:             "/hire",
			expectedResponse: expectedData,
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
			resp, err := http.Post(url, c.contentType, bytes.NewBuffer(jsonData))

			if err != nil {
				t.Error(err)
				return
			}

			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
				return
			}

			actualResponse := string(respBody)
			if !(strings.Contains(actualResponse, c.expectedResponse)) {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, actualResponse)
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
	expectedData := `{"firstName": "John", "lastName": "Doe"}`
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
			expectedResponse: expectedData,
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
			event.SetData(cloudevents.ApplicationJSON, c.data)

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

			actualResponse := string(reqEvent.Data())
			if !(strings.Contains(actualResponse, c.expectedResponse)) {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, actualResponse)
			}
		})
	}
}

func TestJavaCloudEventsOverHTTP(t *testing.T) {
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
	expectedData := `{"firstName":"John", "lastName":"Doe"}`
	cases := []struct {
		name string
		tag  string

		methodType       string
		contentType      string
		path             string
		data             []byte
		expectedResponse string
	}{
		{
			name: "Java CloudEvents Gradle",
			tag:  "java-cloudevents-gradle",

			methodType:       http.MethodPost,
			contentType:      "application/json",
			path:             "/hire",
			data:             jsonData,
			expectedResponse: expectedData,
		},
		{
			name: "Java CloudEvents Maven",
			tag:  "java-cloudevents-maven",

			methodType:       http.MethodPost,
			contentType:      "application/json",
			path:             "/hire",
			data:             jsonData,
			expectedResponse: expectedData,
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

			resp, err := http.Post(url, c.contentType, bytes.NewBuffer(jsonData))

			if err != nil {
				t.Error(err)
				return
			}

			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
				return
			}

			actualResponse := string(respBody)
			if !(strings.Contains(actualResponse, c.expectedResponse)) {
				t.Errorf("Expected response '%s' but received '%s'.", c.expectedResponse, actualResponse)
			}
		})
	}
}
