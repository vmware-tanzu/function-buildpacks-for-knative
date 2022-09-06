RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
buildpacks.python.path := $(abspath $(path))

include $(buildpacks.python.path)/../rules.mk

$(eval $(call VERSION_template,buildpacks.python,$(buildpacks.python.path)))
buildpacks.python.version := $(buildpacks.python.version.release)
buildpacks.python.build_dir := $(buildpack.build_dir)/python/$(buildpacks.python.version)
buildpacks.python.temp_dir := $(buildpack.build_dir)/python/temp/$(buildpacks.python.version)
buildpacks.python.out_dir := $(buildpack.out_dir)/python/$(buildpacks.python.version)

###########################
# Detect and Build binaries
###########################

buildpacks.python.bin.sources := \
	$(buildpacks.python.path)/go.mod \
	$(buildpacks.python.path)/go.sum \
	$(shell find '$(buildpacks.python.path)' -type f -iname '*.go' -not -path '$(buildpacks.python.path)/tests/*') \
	$(shell find '$(buildpacks.python.path)' -type d -not -path '$(buildpacks.python.path)/tests' -mindepth 1) \

buildpacks.python.bin.build := $(buildpacks.python.build_dir)/bin/build
buildpacks.python.bin.detect := $(buildpacks.python.build_dir)/bin/detect
buildpacks.python.bin.main := $(abspath $(buildpacks.python.build_dir)/bin/main)
buildpacks.python.bin := $(buildpacks.python.bin.build) $(buildpacks.python.bin.detect)

$(buildpacks.python.bin.main): $(buildpacks.python.bin.sources)
	@mkdir -p $(@D)
	cd $(buildpacks.python.path) && GOOS="linux" go build -ldflags='-s -w' -o $@ kn-fn/python-function-buildpack/cmd/main

$(buildpacks.python.bin): $(buildpacks.python.bin.main)
	@mkdir -p $(@D)
	ln -sf $(notdir $(buildpacks.python.bin.main)) $@

##################
# Invoker Metadata
##################

buildpacks.python.invoker.version := $(shell cat $(buildpacks.python.path)/../../invokers/python/VERSION)
buildpacks.python.invoker.url := $(base_url)/releases/download/python-invoker/v$(buildpacks.python.invoker.version)/pyfunc-invoker-$(buildpacks.python.invoker.version).tar.gz
buildpacks.python.invoker.sha := $(buildpacks.python.temp_dir)/invoker.sha256
buildpacks.python.invoker_deps.url := $(base_url)/releases/download/python-invoker/v$(buildpacks.python.invoker.version)/pyfunc-invoker-deps-$(buildpacks.python.invoker.version).tar.gz
buildpacks.python.invoker_deps.sha := $(buildpacks.python.temp_dir)/invoker_dep.sha256

$(buildpacks.python.invoker.sha):
	@mkdir -p $(@D)
	curl -sw '\0' -L $(buildpacks.python.invoker.url).sha256 -o $@

$(buildpacks.python.invoker_deps.sha):
	@mkdir -p $(@D)
	curl -sw '\0' -L $(buildpacks.python.invoker_deps.url).sha256 -o $@

buildpacks.python.invoker.values := $(buildpacks.python.temp_dir)/values.yaml

# If modifying this content make sure you're using SPACES and not TABS for indentation
define buildpacks.python.values.content
#@data/values
---
buildpack:
  version: $(buildpacks.python.version)

invoker:
  version: $(buildpacks.python.invoker.version)
  url: $(buildpacks.python.invoker.url)
  sha: $(firstword $(shell cat $(buildpacks.python.invoker.sha)))

invoker_dep:
  version: $(buildpacks.python.invoker.version)
  url: $(buildpacks.python.invoker_deps.url)
  sha: $(firstword $(shell cat $(buildpacks.python.invoker_deps.sha)))
endef

$(buildpacks.python.invoker.values): $(buildpacks.python.invoker.sha) $(buildpacks.python.invoker_deps.sha)
	@mkdir -p $(@D)
	printf '$(subst $(newline),\n,$(buildpacks.python.values.content))' > $@

################
# Buildpack.toml
################

buildpacks.python.config := $(buildpacks.python.build_dir)/buildpack.toml
buildpacks.python.config_intermediate := $(buildpacks.python.temp_dir)/buildpack.yaml

buildpacks.python.config.sources := \
	$(buildpacks.python.path)/ytt \
	$(shell find '$(buildpacks.python.path)/ytt' -type f) \

.INTERMEDIATE: $(buildpacks.python.config_intermediate)
$(buildpacks.python.config_intermediate): $(YTT) $(buildpacks.python.config.sources) $(buildpacks.python.invoker.values)
	@mkdir -p $(@D)
	$(YTT) -f $(buildpacks.python.path)/ytt -f $(buildpacks.python.invoker.values) --ignore-unknown-comments > $@

$(buildpacks.python.config): $(YJ) $(buildpacks.python.config_intermediate)
	@mkdir -p $(@D)
	$(YJ) -yt < $(buildpacks.python.config_intermediate) > $@

#####################
# Buildpack Packaging
# Online = Buildpack does NOT contain the invoker dependencies (NOT cached)
# Offline = Buildpack contains the invoker dependencies (cached)
#####################

buildpacks.python.package.online := $(buildpacks.python.temp_dir)/packages/online
buildpacks.python.package.offline := $(buildpacks.python.temp_dir)/packages/offline

$(buildpacks.python.package.online)/package.toml $(buildpacks.python.package.offline)/package.toml: %/package.toml:
	@mkdir -p $(@D)
	printf '[buildpack]\nuri = "%s"\n\n[platform]\nos = "%s"\n' $* "linux" > $@

buildpacks.python.package.deps := \
	$(buildpacks.python.build_dir)/README.md \
	$(buildpacks.python.bin) \
	$(buildpacks.python.main) \
	$(buildpacks.python.config) \
	$(buildpacks.python.build_dir)/VERSION \

$(buildpacks.python.build_dir)/%: $(buildpacks.python.path)/%
	@mkdir -p $(@D)
	cp $< $@

$(buildpacks.python.package.offline): buildpacks.python.package.flags=--include-dependencies
$(buildpacks.python.package.online) $(buildpacks.python.package.offline): $(CREATE-PACKAGE) $(buildpacks.python.package.deps)
	@mkdir -p $(@D)
	$(CREATE-PACKAGE) \
		$(buildpacks.python.package.flags) \
		--cache-location "$(buildpacks.python.temp_dir)/cache/carton-cache" \
		--destination "$@" \
		--source "$(buildpacks.python.build_dir)"

buildpacks.python.images.online.name := $(registry)/python-buildpack:$(buildpacks.python.version)
buildpacks.python.images.online.out := $(buildpacks.python.out_dir)/python-buildpack-$(buildpacks.python.version).image
buildpacks.python.images.offline.name := $(registry)/python-buildpack-with-deps:$(buildpacks.python.version)
buildpacks.python.images.offline.out := $(buildpacks.python.out_dir)/python-buildpack-with-deps-$(buildpacks.python.version).image

$(buildpacks.python.images.online.out): $(PACK) $(buildpacks.python.package.online) $(buildpacks.python.package.online)/package.toml
	@mkdir -p $(@D)
	$(PACK) buildpack package $(buildpacks.python.images.online.name) --config $(buildpacks.python.package.online)/package.toml
	printf '$(buildpacks.python.images.online.name)' > $@

$(buildpacks.python.images.offline.out): $(PACK) $(buildpacks.python.package.offline) $(buildpacks.python.package.offline)/package.toml
	@mkdir -p $(@D)
	$(PACK) buildpack package $(buildpacks.python.images.offline.name) --config $(buildpacks.python.package.offline)/package.toml
	printf '$(buildpacks.python.images.offline.name)' > $@

########
# commit
########

$(buildpacks.python.path)/buildpack.toml: $(buildpacks.python.config)
	sed 's/^/# /' $(ROOT_DIR)hack/boilerplate.go.txt > $@
	printf '#\n# DO NOT EDIT!! This file is auto-generated.\n\n' >> $@
	cat $< >> $@

#########
# Testing
#########

buildpacks.python.tests.sources := \
	$(buildpacks.python.path)/go.mod \
	$(buildpacks.python.path)/go.sum \
	$(shell find '$(buildpacks.python.path)/tests' -type f -iname '*.go') \
	$(buildpacks.python.path)/tests \

#########
# Targets
#########

buildpacks.python.images.offline: $(buildpacks.python.images.offline.out) $(buildpacks.python.images.offline.out).sha256
buildpacks.python.images.online: $(buildpacks.python.images.online.out) $(buildpacks.python.images.online.out).sha256
buildpacks.python.images: buildpacks.python.images.offline buildpacks.python.images.online

buildpacks.python.images.offline.publish: buildpacks.python.images.offline
	docker push $(buildpacks.python.images.offline.name)
buildpacks.python.images.online.publish: buildpacks.python.images.online
	docker push $(buildpacks.python.images.online.name)
buildpacks.python.images.publish: buildpacks.python.images.offline.publish buildpacks.python.images.online.publish

buildpacks.python.commit.prep: $(buildpacks.python.path)/buildpack.toml

buildpacks.python.tests: $(buildpacks.python.bin.sources) $(buildpacks.python.tests.sources)
	cd $(buildpacks.python.path) && go test -v -count=1 -timeout 30s kn-fn/python-function-buildpack/tests

buildpacks.python.clean:
	rm -rf $(buildpacks.python.build_dir)
	rm -rf $(buildpacks.python.temp_dir)
	rm -rf $(buildpacks.python.out_dir)

buildpacks.offline .PHONY: buildpacks.python.images.offline
buildpacks.offline.publish .PHONY: buildpacks.python.images.offline.publish
buildpacks.online .PHONY: buildpacks.python.images.online
buildpacks.online.publish .PHONY: buildpacks.python.images.online.publish
buildpacks.publish .PHONY: buildpacks.offline.publish buildpacks.online.publish
buildpacks.tests .PHONY: buildpacks.python.tests
buildpacks.clean .PHONY: buildpacks.python.clean
