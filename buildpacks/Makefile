RULES.MK ?= ../rules.mk
include $(RULES.MK)

ifndef buildpack.path

path ?= .
buildpack.path := $(abspath $(path))

include $(buildpack.path)/rules.mk

# Include test case images
buildpack.subdirs = \
	common \
	python \
	java \

SUBDIRS := $(addprefix $(buildpack.path)/, $(buildpack.subdirs))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))
endif

buildpacks.offline:
buildpacks .PHONY: buildpacks.offline

buildpacks.online:
buildpacks .PHONY: buildpacks.online

all .PHONY: buildpacks

buildpacks.clean:
clean .PHONY: buildpacks.clean

buildpacks.tests:
tests .PHONY: buildpacks.tests
