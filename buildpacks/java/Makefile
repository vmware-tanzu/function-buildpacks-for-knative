RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
buildpacks.java.path := $(abspath $(path))

include $(buildpacks.java.path)/../rules.mk

$(eval $(call VERSION_template,buildpacks.java,$(buildpacks.java.path)))
buildpacks.java.version := $(buildpacks.java.version.release)
buildpacks.java.build_dir := $(buildpack.build_dir)/java/$(buildpacks.java.version)
buildpacks.java.temp_dir := $(buildpack.build_dir)/java/temp/$(buildpacks.java.version)
buildpacks.java.out_dir := $(buildpack.out_dir)/java/$(buildpacks.java.version)

###########################
# Detect and Build binaries
###########################

buildpacks.java.bin.sources := \
	$(buildpacks.java.path)/go.mod \
	$(buildpacks.java.path)/go.sum \
	$(shell find '$(buildpacks.java.path)' -type f -iname '*.go' -not -path '$(buildpacks.java.path)/tests/*') \
	$(shell find '$(buildpacks.java.path)' -type d -not -path '$(buildpacks.java.path)/tests' -mindepth 1) \

buildpacks.java.bin.build := $(buildpacks.java.build_dir)/bin/build
buildpacks.java.bin.detect := $(buildpacks.java.build_dir)/bin/detect
buildpacks.java.bin.main := $(abspath $(buildpacks.java.build_dir)/bin/main)
buildpacks.java.bin := $(buildpacks.java.bin.build) $(buildpacks.java.bin.detect)

$(buildpacks.java.bin.main): $(buildpacks.java.bin.sources)
	@mkdir -p $(@D)
	cd $(buildpacks.java.path) && GOOS="linux" go build -ldflags='-s -w' -o $@ kn-fn/java-function-buildpack/cmd/main

$(buildpacks.java.bin): $(buildpacks.java.bin.main)
	@mkdir -p $(@D)
	ln -sf $(notdir $(buildpacks.java.bin.main)) $@

##################
# Invoker Metadata
##################

buildpacks.java.invoker.version := $(shell cat $(buildpacks.java.path)/../../invokers/java/VERSION)
buildpacks.java.invoker.url := $(base_url)/releases/download/java-invoker/v$(buildpacks.java.invoker.version)/java-function-invoker-$(buildpacks.java.invoker.version).jar
buildpacks.java.invoker.sha := $(buildpacks.java.temp_dir)/invoker.sha256

$(buildpacks.java.invoker.sha):
	@mkdir -p $(@D)
	curl -sw '\0' -L $(buildpacks.java.invoker.url).sha256 -o $@

buildpacks.java.invoker.values := $(buildpacks.java.temp_dir)/values.yaml

# If modifying this content make sure you're using SPACES and not TABS for indentation
define buildpacks.java.values.content
#@data/values
---
buildpack:
  version: $(buildpacks.java.version)

invoker:
  version: $(buildpacks.java.invoker.version)
  url: $(buildpacks.java.invoker.url)
  sha: $(firstword $(shell cat $(buildpacks.java.invoker.sha)))
endef

$(buildpacks.java.invoker.values): $(buildpacks.java.invoker.sha)
	@mkdir -p $(@D)
	printf '$(subst $(newline),\n,$(buildpacks.java.values.content))' > $@

################
# Buildpack.toml
################

buildpacks.java.config := $(buildpacks.java.build_dir)/buildpack.toml
buildpacks.java.config_intermediate := $(buildpacks.java.temp_dir)/buildpack.yaml

buildpacks.java.config.sources := \
	$(buildpacks.java.path)/ytt \
	$(shell find '$(buildpacks.java.path)/ytt' -type f) \

.INTERMEDIATE: $(buildpacks.java.config_intermediate)
$(buildpacks.java.config_intermediate): $(YTT) $(buildpacks.java.config.sources) $(buildpacks.java.invoker.values)
	@mkdir -p $(@D)
	$(YTT) -f $(buildpacks.java.path)/ytt -f $(buildpacks.java.invoker.values) --ignore-unknown-comments > $@

$(buildpacks.java.config): $(YJ) $(buildpacks.java.config_intermediate)
	@mkdir -p $(@D)
	$(YJ) -yt < $(buildpacks.java.config_intermediate) > $@

#####################
# Buildpack Packaging
# Online = Buildpack does NOT contain the invoker dependencies (NOT cached)
# Offline = Buildpack contains the invoker dependencies (cached)
#####################

buildpacks.java.package.online := $(buildpacks.java.temp_dir)/packages/online
buildpacks.java.package.offline := $(buildpacks.java.temp_dir)/packages/offline

$(buildpacks.java.package.online)/package.toml $(buildpacks.java.package.offline)/package.toml: %/package.toml:
	@mkdir -p $(@D)
	printf '[buildpack]\nuri = "%s"\n\n[platform]\nos = "%s"\n' $* "linux" > $@

buildpacks.java.package.deps := \
	$(buildpacks.java.build_dir)/README.md \
	$(buildpacks.java.bin) \
	$(buildpacks.java.main) \
	$(buildpacks.java.config) \
	$(buildpacks.java.build_dir)/VERSION \

$(buildpacks.java.build_dir)/%: $(buildpacks.java.path)/%
	@mkdir -p $(@D)
	cp $< $@

$(buildpacks.java.package.offline): buildpacks.java.package.flags=--include-dependencies
$(buildpacks.java.package.online) $(buildpacks.java.package.offline): $(CREATE-PACKAGE) $(buildpacks.java.package.deps)
	@mkdir -p $(@D)
	$(CREATE-PACKAGE) \
		$(buildpacks.java.package.flags) \
		--cache-location "$(buildpacks.java.temp_dir)/cache/carton-cache" \
		--destination "$@" \
		--source "$(buildpacks.java.build_dir)"

buildpacks.java.images.online.name := $(registry)/java-buildpack:$(buildpacks.java.version)
buildpacks.java.images.online.out := $(buildpacks.java.out_dir)/java-buildpack-$(buildpacks.java.version).image
buildpacks.java.images.offline.name := $(registry)/java-buildpack-with-deps:$(buildpacks.java.version)
buildpacks.java.images.offline.out := $(buildpacks.java.out_dir)/java-buildpack-with-deps-$(buildpacks.java.version).image

$(buildpacks.java.images.online.out): $(PACK) $(buildpacks.java.package.online) $(buildpacks.java.package.online)/package.toml
	@mkdir -p $(@D)
	$(PACK) buildpack package $(buildpacks.java.images.online.name) --config $(buildpacks.java.package.online)/package.toml
	printf '$(buildpacks.java.images.online.name)' > $@

$(buildpacks.java.images.offline.out): $(PACK) $(buildpacks.java.package.offline) $(buildpacks.java.package.offline)/package.toml
	@mkdir -p $(@D)
	$(PACK) buildpack package $(buildpacks.java.images.offline.name) --config $(buildpacks.java.package.offline)/package.toml
	printf '$(buildpacks.java.images.offline.name)' > $@

########
# commit
########

$(buildpacks.java.path)/buildpack.toml: $(buildpacks.java.config)
	sed 's/^/# /' $(ROOT_DIR)hack/boilerplate.go.txt > $@
	printf '#\n# DO NOT EDIT!! This file is auto-generated.\n\n' >> $@
	cat $< >> $@

#########
# Testing
#########

buildpacks.java.tests.sources := \
	$(buildpacks.java.path)/go.mod \
	$(buildpacks.java.path)/go.sum \
	$(shell find '$(buildpacks.java.path)/tests' -type f -iname '*.go') \
	$(buildpacks.java.path)/tests \

#########
# Targets
#########

buildpacks.java.images.local: buildpacks.java.invoker.url=file://$(abspath $(invokers.out_dir)/java/java-function-invoker-$(buildpacks.java.invoker.version).jar)
buildpacks.java.images.local: buildpacks.java.images.offline

buildpacks.java.images.offline: $(buildpacks.java.images.offline.out) $(buildpacks.java.images.offline.out).sha256
buildpacks.java.images.online: $(buildpacks.java.images.online.out) $(buildpacks.java.images.online.out).sha256
buildpacks.java.images: buildpacks.java.images.offline buildpacks.java.images.online

buildpacks.java.images.offline.publish: buildpacks.java.images.offline
	docker push $(buildpacks.java.images.offline.name)
buildpacks.java.images.online.publish: buildpacks.java.images.online
	docker push $(buildpacks.java.images.online.name)
buildpacks.java.images.publish: buildpacks.java.images.offline.publish buildpacks.java.images.online.publish

buildpacks.java.commit.prep: $(buildpacks.java.path)/buildpack.toml

buildpacks.java.tests: $(buildpacks.java.bin.sources) $(buildpacks.java.tests.sources)
	cd $(buildpacks.java.path) && go test -v -count=1 -timeout 30s kn-fn/java-function-buildpack/tests

buildpacks.java.clean:
	rm -rf $(buildpacks.java.path)/out
	rm -rf $(buildpacks.java.path)/bin
	rm -rf $(buildpacks.java.path)/package.toml

all buildpacks .PHONY: buildpacks.java.images
buildpacks.offline .PHONY: buildpacks.java.images.offline
buildpacks.offline.publish .PHONY: buildpacks.java.images.offline.publish
buildpacks.online .PHONY: buildpacks.java.images.online
buildpacks.online.publish .PHONY: buildpacks.java.images.online.publish
buildpacks.publish .PHONY: buildpacks.offline.publish buildpacks.online.publish
tests buildpacks.tests .PHONY: buildpacks.java.tests
clean .PHONY: buildpacks.java.clean
