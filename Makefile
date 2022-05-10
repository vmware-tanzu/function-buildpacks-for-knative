ADDLICENSE ?= go run github.com/google/addlicense@latest
RULES.MK := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))/rules.mk
include $(RULES.MK)

# Include subdirs
SUBDIRS := invokers buildpacks builder tests
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))

all:

builder:
buildpacks:
buildpack-files:
invokers:

publish-buildpacks:
publish-invokers:
publish-builder:

tests:
buildpack-tests:
invoker-tests:
smoke-tests:

clean:

.PHONY: add-copyright
add-copyright:
	$(ADDLICENSE) -f hack/boilerplate.go.txt .

.PHONY: check-copyright
check-copyright:
	$(ADDLICENSE) -ignore .github/ -f hack/boilerplate.go.txt -check .
