ADDLICENSE ?= go run github.com/google/addlicense@latest
RULES.MK := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))/rules.mk
include $(RULES.MK)

# Include subdirs
SUBDIRS := invokers buildpacks builder tests
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))

all:

builder:
buildpacks:
invokers:

buildpacks.publish:
invokers.publish:
builder.publish:

tests:
buildpacks.tests:
invokers.tests:
smoke-tests:

clean:

.PHONY: add-copyright
add-copyright:
	$(ADDLICENSE) -f hack/boilerplate.go.txt .

.PHONY: check-copyright
check-copyright:
	$(ADDLICENSE) -ignore "**/*.json" -ignore ".github/**" -ignore "**/config/*.yaml" -f hack/boilerplate.go.txt -check .
