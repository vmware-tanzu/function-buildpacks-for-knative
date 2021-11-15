RULES.MK := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))/rules.mk
include $(RULES.MK)

# Include subdirs
SUBDIRS := invokers buildpacks builder tests
$(foreach dir,$(SUBDIRS),$(eval $(INCLUDE_FILE)))

all:
publish:
builder:
buildpacks:
buildpack-files:
invokers:
publish-buildpacks:
publish-invokers:
publish-builder:
tests:
clean:
