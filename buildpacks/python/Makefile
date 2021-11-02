DETECT_BIN := bin/detect
BUILD_BIN := bin/build
MAIN_BIN := bin/main

SOURCES := $(shell find . -type f -iname '*.go')

.PHONY: all
build: $(BUILD_BIN) $(DETECT_BIN)

$(BUILD_BIN) $(DETECT_BIN): $(MAIN_BIN)
	mkdir -p ${@D}
	ln -sf $(notdir $(MAIN_BIN)) $@

$(MAIN_BIN): $(SOURCES)
	mkdir -p ${@D}
	GOOS="linux" go build -ldflags='-s -w' -o $@ kn-fn/python-function-buildpack/cmd/main

.PHONY: tests
tests:
	go test -count=1 -timeout 30s kn-fn/python-function-buildpack/tests

BUILDPACK_NAME ?= kn-fn/python-function:$(IMAGE_TAG)
BUILDPACK_FILE := out/package.$(IMAGE_TAG).cnb
BUILDPACK_IMAGE := $(IMAGE_REGISTRY)/python-buildpack:$(IMAGE_TAG)

.PHONY: package
package: $(PACK) $(BUILD_BIN) $(DETECT_BIN) buildpack.toml package.toml
	$(PACK) buildpack package $(BUILDPACK_NAME) --config package.toml

.PHONY: package-file
package-file: $(PACK) $(BUILDPACK_FILE)
$(BUILDPACK_FILE):
	mkdir -p ${@D}
	$(PACK) buildpack package $(BUILDPACK_FILE) --config package.toml --format file

.PHONY: publish
publish: package
	docker tag $(BUILDPACK_NAME) $(BUILDPACK_IMAGE)
	docker push $(BUILDPACK_IMAGE)

.PHONY: clean
clean:
	rm -rf out/
	rm -rf bin/
