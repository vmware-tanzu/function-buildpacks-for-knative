RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
tests.path := $(abspath $(path))

# Include test case images
SUBDIRS := $(wildcard $(addprefix $(tests.path), /testdata/template_and_smoke/*))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))

smoke-tests.run: smoke-tests.images
	cd $(tests.path) && go test -v -tags=smoke -count=1

tests smoke-tests .PHONY: smoke-tests.run

template-tests.run: template-tests.images
	cd $(tests.path) && go test -v -tags=template -count=1

tests template-tests .PHONY: template-tests.run

integration-tests.run: $(PACK)
	cd $(tests.path) && go test ./integration/... -v -count=1 -args -pack=$(PACK) -builder=paketobuildpacks/builder:base

tests integration-tests .PHONY: integration-tests.run
