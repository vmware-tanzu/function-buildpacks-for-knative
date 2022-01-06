RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
python-buildpack.path := $(abspath $(path))

$(eval $(call VERSION_template,python-buildpack,$(python-buildpack.path)))

python-buildpack.go-sources := \
	$(python-buildpack.path)/go.mod \
	$(python-buildpack.path)/go.sum \
	$(shell find '$(python-buildpack.path)' -type f -iname '*.go' -not -path '$(python-buildpack.path)/tests/*') \
	$(shell find '$(python-buildpack.path)' -type d -not -path '$(python-buildpack.path)/tests' -mindepth 1) \

python-buildpack.test-sources := \
	$(python-buildpack.path)/go.mod \
	$(python-buildpack.path)/go.sum \
	$(shell find '$(python-buildpack.path)/tests' -type f -iname '*.go') \
	$(python-buildpack.path)/tests \

python-buildpack.config-sources := $(python-buildpack.path)/buildpack.toml $(python-buildpack.path)/package.toml
python-buildpack.image := $(registry)/python-buildpack:$(python-buildpack.version)
python-buildpack.package-file := $(python-buildpack.path)/out/package.$(python-buildpack.version).cnb

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

python-buildpack.tests: $(python-buildpack.go-sources) $(python-buildpack.test-sources)
	cd $(python-buildpack.path) && go test -v -count=1 -timeout 30s kn-fn/python-function-buildpack/tests

python-buildpack.clean:
	rm -rf $(python-buildpack.path)/out
	rm -rf $(python-buildpack.path)/bin

all buildpacks .PHONY: python-buildpack.package
buildpack-files: $(python-buildpack.package-file)
publish-buildpacks .PHONY: python-buildpack.publish
tests buildpack-tests .PHONY: python-buildpack.tests
clean .PHONY: python-buildpack.clean
