ADDLICENSE ?= go run github.com/google/addlicense@latest
RULES.MK := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))/rules.mk
include $(RULES.MK)

# Include subdirs
SUBDIRS := invokers buildpacks tests
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))

all:

buildpacks:
invokers:

buildpacks.publish:
invokers.publish:

tests:
buildpacks.tests:
invokers.tests:
smoke-tests:
template-tests:
integration-tests:

.PHONY: clean main.clean
main.clean:
	$(RM) -rf $(out_dir)
	$(RM) -rf $(build_dir)
clean: main.clean

.PHONY: add-copyright
add-copyright:
	$(ADDLICENSE) -f hack/boilerplate.go.txt .

.PHONY: check-copyright
check-copyright:
	$(ADDLICENSE) -ignore "**/*.json" -ignore ".github/**" -ignore "**/config/*.yaml" -f hack/boilerplate.go.txt -check .
