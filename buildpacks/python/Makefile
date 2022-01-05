RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
python-buildpack.path := $(abspath $(path))

python-buildpack.go-sources := $(shell find $(python-buildpack.path) -type f -iname '*.go')
python-buildpack.config-sources := $(python-buildpack.path)/buildpack.toml $(python-buildpack.path)/package.toml
python-buildpack.image := $(IMAGE_REGISTRY)/python-buildpack:$(IMAGE_TAG)
python-buildpack.package-file := $(python-buildpack.path)/out/package.$(IMAGE_TAG).cnb

python-buildpack.build-bin := $(python-buildpack.path)/bin/build
python-buildpack.detect-bin := $(python-buildpack.path)/bin/detect
python-buildpack.main-bin := $(abspath $(python-buildpack.path)/bin/main)
python-buildpack.bin := $(python-buildpack.build-bin) $(python-buildpack.detect-bin)

python-buildpack.deps := $(python-buildpack.bin) $(python-buildpack.config-sources)

$(python-buildpack.build-bin) $(python-buildpack.detect-bin): $(python-buildpack.main-bin)
	@mkdir -p $(@D)
	ln -sf $(notdir $(python-buildpack.main-bin)) $@

$(python-buildpack.main-bin): $(python-buildpack.go-sources)
	@mkdir -p $(@D)
	cd $(python-buildpack.path) && GOOS="linux" go build -ldflags='-s -w' -o $@ kn-fn/python-function-buildpack/cmd/main

python-buildpack.package-deps := $(PACK) $(python-buildpack.build-bin) $(python-buildpack.detect-bin) $(python-buildpack.config-sources)
python-buildpack.package: $(python-buildpack.package-deps)
	$(PACK) buildpack package $(python-buildpack.image) --config $(python-buildpack.path)/package.toml

$(python-buildpack.package-file): $(PACK) $(python-buildpack.package-deps)
	@mkdir -p $(@D)
	$(PACK) buildpack package $@ --config $(python-buildpack.path)/package.toml --format file

python-buildpack.publish: python-buildpack.package
	docker push $(python-buildpack.image)

python-buildpack.tests: $(python-buildpack.go-sources)
	cd $(python-buildpack.path) && go test -count=1 -timeout 30s kn-fn/python-function-buildpack/tests

python-buildpack.clean:
	rm -rf $(python-buildpack.path)/out
	rm -rf $(python-buildpack.path)/bin

all buildpacks .PHONY: python-buildpack.package
buildpack-files: $(python-buildpack.package-file)
publish publish-buildpacks .PHONY: python-buildpack.publish
tests buildpack-tests .PHONY: python-buildpack.tests
clean .PHONY: python-buildpack.clean
