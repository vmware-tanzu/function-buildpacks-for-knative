ifndef RULES_MK # Prevent repeated "-include".
RULES_MK := $(lastword $(MAKEFILE_LIST))
RULES_INCLUDE_DIR := $(dir $(RULES_MK))
ROOT_DIR := $(RULES_INCLUDE_DIR)

-include $(ROOT_DIR)/local.mk

.DEFAULT_GOAL := all
.DELETE_ON_ERROR: # This will delete files from targets that don't succeed.
.SUFFIXES: # This removes a lot of the implicit rules.

os.name := $(shell uname -s | tr A-Z a-z)

out_dir := $(abspath $(ROOT_DIR)/out)
build_dir := $(abspath $(ROOT_DIR)/build)

time_format := epoch
build_time.human := $(shell date +'%Y%m%d-%H%M%S')
build_time.epoch := $(shell date +'%s')
build_time = $(build_time.$(time_format))

git.dirty := $(shell git status -s)
git.branch := $(shell git rev-parse --abbrev-ref HEAD)
git.commit := $(shell git rev-parse HEAD)

base_url := https://github.com/vmware-tanzu/function-buildpacks-for-knative

registry.location ?= ghcr
registry.ghcr := ghcr.io/vmware-tanzu/function-buildpacks-for-knative
registry.other := $(REGISTRY)
registry = $(registry.$(registry.location))

ifeq ($(registry.location), other)
ifndef REGISTRY
$(error REGISTRY not defined. This is required for targeting "other" registry)
endif
endif

build := commit
ifneq (,$(findstring release/, $(git.branch)))
build := release
else ifeq ($(GIT_BRANCH), main)
build := latest
endif

# If the repo is dirty, we will consider it dev.
# Commit your work if you want consistency
ifeq ($(GIT_DIRTY),)
build := dev
endif

include $(ROOT_DIR)/version.mk
include $(ROOT_DIR)/tools.mk

define newline


endef

# Sha recipes
%.sha256: %
	cd $(dir $<) && shasum -a 256 $(notdir $<) > $@

%.print-sha: %.sha256
	@cat $<

define INCLUDE_FILE
path = $(1)
include $(1)/Makefile
endef

endif
