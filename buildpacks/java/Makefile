RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
java-buildpack.path := $(abspath $(path))

$(eval $(call VERSION_template,java-buildpack,$(java-buildpack.path)))

java-buildpack.go-sources := \
	$(java-buildpack.path)/go.mod \
	$(java-buildpack.path)/go.sum \
	$(shell find '$(java-buildpack.path)' -type f -iname '*.go' -not -path '$(java-buildpack.path)/tests/*') \
	$(shell find '$(java-buildpack.path)' -type d -not -path '$(java-buildpack.path)/tests' -mindepth 1) \

java-buildpack.test-sources := \
	$(java-buildpack.path)/go.mod \
	$(java-buildpack.path)/go.sum \
	$(shell find '$(java-buildpack.path)/tests' -type f -iname '*.go') \
	$(java-buildpack.path)/tests \

java-buildpack.config-sources := $(java-buildpack.path)/buildpack.toml $(java-buildpack.path)/package.toml
java-buildpack.image := $(registry)/java-buildpack:$(java-buildpack.version)
java-buildpack.package-file := $(java-buildpack.path)/out/package.$(java-buildpack.version).cnb

java-buildpack.build-bin := $(java-buildpack.path)/bin/build
java-buildpack.detect-bin := $(java-buildpack.path)/bin/detect
java-buildpack.main-bin := $(abspath $(java-buildpack.path)/bin/main)
java-buildpack.bin := $(java-buildpack.build-bin) $(java-buildpack.detect-bin)

java-buildpack.deps := $(java-buildpack.bin) $(java-buildpack.config-sources)

$(java-buildpack.build-bin) $(java-buildpack.detect-bin): $(java-buildpack.main-bin)
	@mkdir -p $(@D)
	ln -sf $(notdir $(java-buildpack.main-bin)) $@

$(java-buildpack.main-bin): $(java-buildpack.go-sources)
	@mkdir -p $(@D)
	cd $(java-buildpack.path) && GOOS="linux" go build -ldflags='-s -w' -o $@ kn-fn/java-function-buildpack/cmd/main

java-buildpack.package-deps := $(PACK) $(java-buildpack.build-bin) $(java-buildpack.detect-bin) $(java-buildpack.config-sources)
java-buildpack.package: $(java-buildpack.package-deps)
	$(PACK) buildpack package $(java-buildpack.image) --config $(java-buildpack.path)/package.toml

$(java-buildpack.package-file): $(PACK) $(java-buildpack.package-deps)
	@mkdir -p $(@D)
	$(PACK) buildpack package $@ --config $(java-buildpack.path)/package.toml --format file

java-buildpack.publish: java-buildpack.package
	docker push $(java-buildpack.image)

java-buildpack.tests: $(java-buildpack.go-sources) $(java-buildpack.test-sources)
	cd $(java-buildpack.path) && go test -v -count=1 -timeout 30s kn-fn/java-function-buildpack/tests

java-buildpack.clean:
	rm -rf $(java-buildpack.path)/out
	rm -rf $(java-buildpack.path)/bin

all buildpacks .PHONY: java-buildpack.package
buildpack-files: $(java-buildpack.package-file)
publish-buildpacks .PHONY: java-buildpack.publish
tests buildpack-tests .PHONY: java-buildpack.tests
clean .PHONY: java-buildpack.clean
