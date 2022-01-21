RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
java-buildpack.path := $(abspath $(path))

include $(java-buildpack.path)/../rules.mk

$(eval $(call VERSION_template,java-buildpack,$(java-buildpack.path)))
java-buildpack.version := $(java-buildpack.version.release)

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
java-buildpack.image_with_deps := $(registry)/java-buildpack-with-deps:$(java-buildpack.version)

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

java-buildpack.package-deps := $(PACK) 
java-buildpack.package.online := $(buildpack.build_dir)/packages/java/$(java-buildpack.version)/online
java-buildpack.package.offline := $(buildpack.build_dir)/packages/java/$(java-buildpack.version)/offline
java-buildpack.online.out := $(buildpack.out_dir)/java/online/java-buildpack-$(java-buildpack.version).image
java-buildpack.offline.out := $(buildpack.out_dir)/java/offline/java-buildpack-$(java-buildpack.version).cnb

$(java-buildpack.package.online): $(CREATE-PACKAGE) $(java-buildpack.deps)
	@mkdir -p $(@D)
	$(CREATE-PACKAGE) \
		--cache-location "$(buildpack.build_dir)/cache/carton-cache" \
		--destination "$@" \
		--source "$(java-buildpack.path)" \
		--version $(java-buildpack.version)

java-buildpack.online.build := $(buildpack.build_dir)/java/online/java-buildpack-$(java-buildpack.version).image
$(java-buildpack.online.build): $(java-buildpack.package.online) $(PACK)
	@mkdir -p $(@D)
	cp $(java-buildpack.path)/package.toml $(java-buildpack.package.online)
	$(PACK) buildpack package $(java-buildpack.image) --config $(java-buildpack.package.online)/package.toml
	printf '$(java-buildpack.image)' > $@

$(java-buildpack.online.out): $(java-buildpack.online.build)
	@mkdir -p $(@D)
	cp $< $@

java-buildpack.online: $(java-buildpack.online.out) $(java-buildpack.online.out).sha256

$(java-buildpack.package.offline): $(CREATE-PACKAGE) $(java-buildpack.deps)
	@mkdir -p $(@D)
	$(CREATE-PACKAGE) \
		--cache-location "$(buildpack.build_dir)/cache/carton-cache" \
		--destination "$@" \
		--source "$(java-buildpack.path)" \
		--include-dependencies \
		--version $(java-buildpack.version)

java-buildpack.offline.build := $(buildpack.build_dir)/java/offline/java-buildpack-$(java-buildpack.version).image
$(java-buildpack.offline.build): $(java-buildpack.package.offline) $(PACK)
	@mkdir -p $(@D)
	cp $(java-buildpack.path)/package.toml $(java-buildpack.package.offline)
	$(PACK) buildpack package $(java-buildpack.image_with_deps) --config $(java-buildpack.package.offline)/package.toml
	printf '$(python-buildpack.image_with_deps)' > $@

$(java-buildpack.offline.out): $(java-buildpack.offline.build)
	@mkdir -p $(@D)
	cp $< $@

java-buildpack.offline: $(java-buildpack.offline.out) $(java-buildpack.offline.out).sha256

java-buildpack.package: java-buildpack.online java-buildpack.offline

java-buildpack.publish: java-buildpack.package
	docker push $(java-buildpack.image)
	docker push $(java-buildpack.image_with_deps)

java-buildpack.tests: $(java-buildpack.go-sources) $(java-buildpack.test-sources)
	cd $(java-buildpack.path) && go test -v -count=1 -timeout 30s kn-fn/java-function-buildpack/tests

java-buildpack.clean:
	rm -rf $(java-buildpack.path)/out
	rm -rf $(java-buildpack.path)/bin

all buildpacks .PHONY: java-buildpack.package
buildpacks.offline .PHONY: java-buildpack.offline
buildpacks.online .PHONY: java-buildpack.online
publish-buildpacks .PHONY: java-buildpack.publish
tests buildpack-tests .PHONY: java-buildpack.tests
clean .PHONY: java-buildpack.clean
