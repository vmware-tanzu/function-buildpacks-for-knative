RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
buildpacks.common.path := $(abspath $(path))

include $(buildpacks.common.path)/../rules.mk

#########
# Testing
#########

buildpacks.common.tests.sources := \
	$(buildpacks.common.path)/go.mod \
	$(buildpacks.common.path)/go.sum \
	$(shell find '$(buildpacks.common.path)/tests' -type f -iname '*.go') \
	$(buildpacks.common.path)/tests \

#########
# Targets
#########

buildpacks.common.tests: $(buildpacks.common.tests.sources)
	cd $(buildpacks.common.path) && go test -v -count=1 -timeout 30s -race -covermode=atomic -coverprofile=coverage.out kn-fn/buildpacks/tests

buildpacks.tests .PHONY: buildpacks.common.tests
