// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package integration

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/onsi/gomega/format"
	"github.com/paketo-buildpacks/occam"
	"github.com/paketo-buildpacks/packit/v2/pexec"

	"github.com/onsi/gomega"
)

func init() {
	format.MaxLength = 0
}

func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port), nil
}

func NewDocker() occam.Docker {
	return occam.NewDocker()
}

var PackExecutable = flag.String("pack", "pack", "Pack executable for the test")
var Builder = flag.String("builder", "", "Builder for the test")

type BuildOpts func(occam.PackBuild) occam.PackBuild

func Pack() occam.Pack {
	return occam.NewPack().WithExecutable(pexec.NewExecutable(*PackExecutable)).WithNoColor()
}

func PackBuild(name string, source string, opts ...BuildOpts) (occam.Image, fmt.Stringer, error) {
	b := Pack().Build.
		WithBuilder(*Builder)

	for _, opt := range opts {
		b = opt(b)
	}

	return b.Execute(name, source)
}

func withEnvs(envs map[string]string) BuildOpts {
	return func(b occam.PackBuild) occam.PackBuild {
		return b.WithEnv(envs)
	}
}

func checkRequest(t *testing.T, method string, port string, path string, data string, contentType string) (int, string) {
	expect := gomega.NewWithT(t).Expect

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	var response *http.Response
	var err error
	if contentType == "" {
		contentType = "application/json"
	}

	if method == http.MethodPost {
		response, err = http.Post(fmt.Sprintf("http://localhost:%s%s", port, path), contentType, bytes.NewBufferString(data))
		expect(err).ToNot(gomega.HaveOccurred())
		defer response.Body.Close()

	} else if method == http.MethodGet {
		response, err = http.Get(fmt.Sprintf("http://localhost:%s%s", port, path))
		expect(err).ToNot(gomega.HaveOccurred())
		defer response.Body.Close()
	}

	expect(err).ToNot(gomega.HaveOccurred())

	content, err := io.ReadAll(response.Body)
	expect(err).NotTo(gomega.HaveOccurred())
	return response.StatusCode, string(content)
}
