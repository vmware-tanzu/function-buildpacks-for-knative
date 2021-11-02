BUILDPACKS := \
	python-buildpack \
	java-buildpack \

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(shell cat VERSION)
ifneq (,$(findstring release/, $(GIT_BRANCH)))
IMAGE_TAG := $(VERSION)
else ifeq ($(GIT_BRANCH), master)
IMAGE_TAG := latest
else
IMAGE_TAG := dev-$(subst /,_,$(GIT_BRANCH))
endif
export IMAGE_TAG

IMAGE_REGISTRY := us.gcr.io/daisy-284300/kn-fn
export IMAGE_REGISTRY

### Tools
OS_NAME := $(shell uname -s | tr A-Z a-z)
TOOLS_BIN_DIR := tools/bin

export PACK := $(abspath $(TOOLS_BIN_DIR)/pack)

TOOLS := $(PACK)

export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

$(PACK):
	@mkdir -p $(@D)
ifeq (darwin,$(OS_NAME))
	curl -sL https://github.com/buildpacks/pack/releases/download/v0.21.1/pack-v0.21.1-macos.tgz | tar -xz -C $(@D) && chmod +x $@
else
	curl -sL https://github.com/buildpacks/pack/releases/download/v0.21.1/pack-v0.21.1-linux.tgz | tar -xz -C $(@D) && chmod +x $@
endif

.PHONY: install-tools
install-tools: $(TOOLS)

export BUILDER_IMAGE := $(IMAGE_REGISTRY)/builder:$(IMAGE_TAG)

.PHONY: builder
builder: $(PACK) $(BUILDPACKS) builder/builder.toml
	$(PACK) builder create -c builder/builder.toml $(BUILDER_IMAGE)

.PHONY: publish-builder
publish-builder: builder
	docker push $(BUILDER_IMAGE)

CLEAN_BUILDPACKS := $(addprefix clean-,$(BUILDPACKS))
$(CLEAN_BUILDPACKS): clean-%-buildpack:
	cd buildpacks/$* && $(MAKE) clean

PUBLISH_BUILDPACKS := $(addprefix publish-,$(BUILDPACKS))
$(PUBLISH_BUILDPACKS): publish-%-buildpack: %-buildpack
	cd buildpacks/$* && $(MAKE) publish

$(BUILDPACKS): %-buildpack: $(PACK)
	cd buildpacks/$* && $(MAKE) build

.PHONY: publish
publish: $(PUBLISH_BUILDPACKS) publish-builder

.PHONY: buildpacks
buildpacks: $(BUILDPACKS)

BUILDPACK_TEST := $(addprefix test-,$(BUILDPACKS))
$(BUILDPACK_TEST): test-%-buildpack:
	cd buildpacks/$* && $(MAKE) tests

.PHONY: smoke-test
smoke-test:
	cd tests && go test -v

.PHONY: tests
tests: $(BUILDPACK_TEST)

.PHONY: clean
clean: $(CLEAN_BUILDPACKS)
	-docker rmi -f $(BUILDER_IMAGE)
	rm -rf tools/

.PHONY: test-images
test-images: $(PACK)
	cd tests/cases/helloworld && $(MAKE) images