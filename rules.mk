ifndef RULES_MK # Prevent repeated "-include".
RULES_MK := $(lastword $(MAKEFILE_LIST))
ROOT_DIR := $(dir $(RULES_MK))

.DEFAULT_GOAL := all
.DELETE_ON_ERROR: # This will delete files from targets that don't succeed.
.SUFFIXES: # This removes a lot of the implicit rules.

OUT_DIR := $(abspath $(ROOT_DIR)/out)
BUILD_DIR := $(abspath $(ROOT_DIR)/build)

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git rev-parse HEAD)
VERSION.release := $(shell cat $(ROOT_DIR)/VERSION)
VERSION.dev := dev-$(subst /,_,$(GIT_BRANCH))
VERSION.latest := latest

build := dev
ifneq (,$(findstring release/, $(GIT_BRANCH)))
build := release
else ifeq ($(GIT_BRANCH), main)
build := latest
endif

IMAGE_TAG := $(VERSION.$(build))
IMAGE_REGISTRY := us.gcr.io/daisy-284300/kn-fn
BUILDER_IMAGE := $(IMAGE_REGISTRY)/builder:$(IMAGE_TAG)

OS_NAME := $(shell uname -s | tr A-Z a-z)

# Sha recipes
%.sha256: %
	cd $(dir $<) && shasum -a 256 $(notdir $<) > $@

# Define the tools here
TOOLS_DIR := $(abspath $(ROOT_DIR)/tools)
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/bin)

PACK := $(TOOLS_BIN_DIR)/pack
GSUTIL := $(TOOLS_DIR)/gsutil/gsutil

$(PACK).darwin:
	@mkdir -p $(@D)
	curl -sL https://github.com/buildpacks/pack/releases/download/v0.21.1/pack-v0.21.1-macos.tgz | tar -xz -C $(@D)
	touch $@

$(PACK).linux:
	@mkdir -p $(@D)
	curl -sL https://github.com/buildpacks/pack/releases/download/v0.21.1/pack-v0.21.1-linux.tgz | tar -xz -C $(@D)
	touch $@

$(PACK): $(PACK).$(OS_NAME)
	chmod +x $@
	touch $@

$(GSUTIL):
	@mkdir -p $(@D)
	curl -sL https://storage.googleapis.com/pub/gsutil.tar.gz | tar -xz -C $(TOOLS_DIR)

define INCLUDE_FILE
path = $(dir)
include $(dir)/Makefile
endef

rules.clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(OUT_DIR)
	rm -rf $(dir $(TOOLS_BIN_DIR))

clean .PHONY: rules.clean

endif
