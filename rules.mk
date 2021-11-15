ifndef RULES_MK # Prevent repeated "-include".
RULES_MK := $(lastword $(MAKEFILE_LIST))
ROOT_DIR := $(dir $(RULES_MK))

.DEFAULT_GOAL := all

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
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

# Define the tools here
TOOLS_BIN_DIR := $(abspath $(ROOT_DIR)/tools/bin)

PACK := $(abspath $(TOOLS_BIN_DIR)/pack)
TOOLS := $(PACK)

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

define INCLUDE_FILE
path = $(dir)
include $(dir)/Makefile
endef

OUT_DIR := $(abspath $(ROOT_DIR)/out)
$(OUT_DIR):
	@mkdir -p $(@D)

rules.clean:
	rm -rf $(OUT_DIR)

clean: rules.clean

endif
