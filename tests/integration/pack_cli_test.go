// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	. "github.com/paketo-buildpacks/occam/matchers"
)

var packCliTestData = filepath.Join("..", "testdata", "e2e", "pack_cli")

func TestPython(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		sourcePaths    []string
		envs           map[string]string
		contentType    string
		data           string
		buildSuccess   bool
		expectedResult string
	}{
		{
			name: "HTTP without func.yaml or BP_FUNCTION",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "python", "http_default"),
			},
		},
		{
			name: "HTTP with BP_FUNCTION only",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "python", "http_default"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "func.main",
			},
			buildSuccess:   true,
			expectedResult: "Hello World!",
		},
		{
			name: "HTTP with non-default module and function",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "python", "http_custom"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "mymodule.myfunc",
			},
			buildSuccess:   true,
			expectedResult: "Hello World!",
		},
		{
			name: "CloudEvent without func.yaml or BP_FUNCTION",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "python", "cloudevent_default"),
			},
		},
		{
			name: "CloudEvent with BP_FUNCTION only",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "python", "cloudevent_default"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "func.main",
			},
			data: `{
				"specversion" : "1.0",
				"source": "local",
				"type" : "hello",
				"id" : "A234-1234-1234",
				"datacontenttype" : "text/plain",
				"data" : "Kapow!"
			}`,
			buildSuccess:   true,
			expectedResult: "Kapow!",
		},
		{
			name: "CloudEvent with non-default module and function",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "python", "cloudevent_custom"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "mymodule.myfunc",
			},
			data: `{
				"specversion" : "1.0",
				"source": "local",
				"type" : "hello",
				"id" : "A234-1234-1234",
				"datacontenttype" : "text/plain",
				"data" : "Kapow!"
			}`,
			buildSuccess:   true,
			expectedResult: "Kapow!",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			expect := NewWithT(t).Expect
			eventually := NewWithT(t).Eventually

			imageName, err := occam.RandomName()
			expect(err).ToNot((HaveOccurred()))

			var latestImage occam.Image
			var latestImageLogs fmt.Stringer

			for _, src := range c.sourcePaths {
				image, logs, err := PackBuild(imageName, src, withEnvs(c.envs))
				if !c.buildSuccess {
					expect(err).To(HaveOccurred())
					expect(logs).To(ContainLines(ContainSubstring("No buildpack groups passed detection.")))
					return // There're no image, nothing else to test on it
				} else {
					defer NewDocker().Image.Remove.Execute(image.ID)
					expect(logs).To(ContainLines(ContainSubstring("kn-fn/python-function")))

					latestImage = image
					latestImageLogs = logs
				}
			}

			port, err := GetFreePort()
			expect(err).ToNot(HaveOccurred())

			container, err := NewDocker().Container.Run.
				WithPublish(port + ":8080").
				Execute(latestImage.ID)
			expect(err).NotTo(HaveOccurred(), latestImageLogs.String())
			defer NewDocker().Container.Remove.Execute(container.ID)
			defer NewDocker().Volume.Remove.Execute(occam.CacheVolumeNames(imageName))

			eventually(container, 2*time.Minute, 10*time.Second).Should(BeAvailable())

			code, response := checkRequest(t, http.MethodPost, container.HostPort("8080"), "/", c.data, c.contentType)
			expect(code).To(Equal(http.StatusOK))
			expect(response).To(Equal(c.expectedResult))

			code, _ = checkRequest(t, http.MethodGet, container.HostPort("8080"), "/health/ready", "", "")
			expect(code).To(Equal(http.StatusOK))

			code, _ = checkRequest(t, http.MethodGet, container.HostPort("8080"), "/health/live", "", "")
			expect(code).To(Equal(http.StatusOK))
		})
	}
}

func TestJava(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		sourcePaths    []string
		envs           map[string]string
		contentType    string
		data           string
		buildSuccess   bool
		expectedResult interface{}
	}{
		{
			name: "HTTP without func.yaml or BP_FUNCTION",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "java", "http_default"),
			},
		},
		{
			name: "HTTP with BP_FUNCTION only",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "java", "http_default"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "functions.Handler",
			},
			buildSuccess:   true,
			expectedResult: "Hello World!",
		},
		{
			name: "HTTP with non-default module and function",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "java", "http_custom"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "com.example.functions.HelloWorld",
			},
			buildSuccess:   true,
			expectedResult: "Hello World!",
		},
		{
			name: "CloudEvent without func.yaml or BP_FUNCTION",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "java", "cloudevent_default"),
			},
		},
		{
			name: "CloudEvent with BP_FUNCTION only",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "java", "cloudevent_default"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "functions.Handler",
			},
			contentType: "application/cloudevents+json",
			data: `{
				"specversion" : "1.0",
				"type" : "hire",
				"source" : "https://spring.io/",
				"id" : "A234-1234-1234",
				"datacontenttype" : "application/json",
				"data": {
						"firstName": "John",
						"lastName": "Doe"
				}
			}`,
			buildSuccess: true,
			expectedResult: map[string]interface{}{
				"person": map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
				"id": float64(0),
			},
		},
		{
			name: "CloudEvent with non-default module and function",
			sourcePaths: []string{
				filepath.Join(packCliTestData, "java", "cloudevent_custom"),
			},
			envs: map[string]string{
				"BP_FUNCTION": "com.example.functions.Hire",
			},
			contentType: "application/cloudevents+json",
			data: `{
		    "specversion" : "1.0",
		    "type" : "hire",
		    "source" : "https://spring.io/",
		    "id" : "A234-1234-1234",
		    "datacontenttype" : "application/json",
		    "data": {
		        "firstName": "John",
		        "lastName": "Doe"
		    }
		  }`,
			buildSuccess: true,
			expectedResult: map[string]interface{}{
				"person": map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
				"id": float64(0),
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			expect := NewWithT(t).Expect
			eventually := NewWithT(t).Eventually

			imageName, err := occam.RandomName()
			expect(err).ToNot((HaveOccurred()))

			var latestImage occam.Image
			var latestImageLogs fmt.Stringer

			for _, src := range c.sourcePaths {
				image, logs, err := PackBuild(imageName, src, withEnvs(c.envs))
				if !c.buildSuccess {
					expect(err).To(HaveOccurred(), logs.String())
					expect(logs).To(ContainLines(ContainSubstring("No buildpack groups passed detection.")))
					return // There're no image, nothing else to test on it
				} else {
					defer NewDocker().Image.Remove.Execute(image.ID)
					expect(logs).To(ContainLines(ContainSubstring("kn-fn/java-function")))

					latestImage = image
					latestImageLogs = logs
				}
			}

			port, err := GetFreePort()
			expect(err).ToNot(HaveOccurred())

			container, err := NewDocker().Container.Run.
				WithPublish(port + ":8080").
				Execute(latestImage.ID)
			expect(err).NotTo(HaveOccurred(), latestImageLogs.String())
			defer NewDocker().Container.Remove.Execute(container.ID)
			defer NewDocker().Volume.Remove.Execute(occam.CacheVolumeNames(imageName))

			eventually(container, 2*time.Minute, 10*time.Second).Should(BeAvailable())

			code, response := checkRequest(t, http.MethodPost, container.HostPort("8080"), "/", c.data, c.contentType)
			expect(code).To(Equal(http.StatusOK))
			switch v := c.expectedResult.(type) {
			default:
				fmt.Printf("unexpected type %T", v)
			case map[string]interface{}:
				jsonResponse := map[string]interface{}{}
				err := json.Unmarshal([]byte(response), &jsonResponse)
				expect(err).ToNot(HaveOccurred())
				expect(jsonResponse).To(Equal(c.expectedResult))
			case string:
				expect(response).To(Equal(c.expectedResult))
			}

			code, _ = checkRequest(t, http.MethodGet, container.HostPort("8080"), "/actuator/health", "", "")
			expect(code).To(Equal(http.StatusOK))
		})
	}
}
