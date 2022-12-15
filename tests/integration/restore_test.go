// Copyright 2021-2022 VMware, Inc.
// SPDX-License-Identifier: BSD-2-Clause

package integration

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	. "github.com/paketo-buildpacks/occam/matchers"
)

var restoreTestData = filepath.Join("..", "testdata", "restore")

// The invoker is using Spring 3 which requires Java 17.
var java17 = map[string]string{
	"BP_JVM_VERSION": "17",
}

func TestRestorePython(t *testing.T) {
	t.Parallel()
	Expect := NewWithT(t).Expect
	Eventually := NewWithT(t).Eventually

	name, err := occam.RandomName()
	Expect(err).ToNot((HaveOccurred()))

	source1, err := occam.Source(filepath.Join(restoreTestData, "python", "step1"))
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(source1)

	source2, err := occam.Source(filepath.Join(restoreTestData, "python", "step2"))
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(source2)

	docker := NewDocker()

	// Do the first build
	{
		image, logs, err := PackBuild(name, source1, withEnvs(java17))
		Expect(err).NotTo(HaveOccurred(), logs.String())
		defer docker.Image.Remove.Execute(image.ID)
	}

	// Do the second build
	image, logs, err := PackBuild(name, source2, withEnvs(java17))
	Expect(err).NotTo(HaveOccurred(), logs.String())
	defer docker.Image.Remove.Execute(image.ID)

	port, err := GetFreePort()
	Expect(err).ToNot(HaveOccurred())

	container, err := docker.Container.Run.
		WithPublish(port + ":8080").
		Execute(image.ID)
	Expect(err).NotTo(HaveOccurred(), logs.String())
	defer docker.Container.Remove.Execute(container.ID)
	defer docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))

	Eventually(container, 2*time.Minute, 10*time.Second).Should(BeAvailable())

	response, err := http.Post(fmt.Sprintf("http://localhost:%s", container.HostPort("8080")), "application/json", bytes.NewBufferString(""))
	Expect(err).NotTo(HaveOccurred())
	defer response.Body.Close()
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	content, err := io.ReadAll(response.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(string(content)).To(Equal("Bye world!"))
}

func TestRestoreJava(t *testing.T) {
	t.Parallel()
	Expect := NewWithT(t).Expect
	Eventually := NewWithT(t).Eventually

	name, err := occam.RandomName()
	Expect(err).ToNot((HaveOccurred()))

	source1, err := occam.Source(filepath.Join(restoreTestData, "java", "step1"))
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(source1)

	source2, err := occam.Source(filepath.Join(restoreTestData, "java", "step2"))
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(source2)

	docker := NewDocker()

	// Do the first build
	{
		image, logs, err := PackBuild(name, source1, withEnvs(java17))
		Expect(err).NotTo(HaveOccurred(), logs.String())
		defer docker.Image.Remove.Execute(image.ID)
	}

	// Do the second build
	image, logs, err := PackBuild(name, source2, withEnvs(java17))
	Expect(err).NotTo(HaveOccurred(), logs.String())
	defer docker.Image.Remove.Execute(image.ID)

	port, err := GetFreePort()
	Expect(err).ToNot(HaveOccurred())

	container, err := docker.Container.Run.
		WithPublish(port + ":8080").
		Execute(image.ID)
	Expect(err).NotTo(HaveOccurred(), logs.String())
	defer docker.Container.Remove.Execute(container.ID)
	defer docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))

	Eventually(container, 2*time.Minute, 10*time.Second).Should(BeAvailable())

	response, err := http.Post(fmt.Sprintf("http://localhost:%s", container.HostPort("8080")), "application/json", bytes.NewBufferString(""))
	Expect(err).NotTo(HaveOccurred())
	defer response.Body.Close()
	Expect(response.StatusCode).To(Equal(http.StatusOK))

	content, err := io.ReadAll(response.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(string(content)).To(Equal("Bye world!"))
}
