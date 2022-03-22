RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
smoke-tests.path := $(abspath $(path))

# Include test case images
SUBDIRS := $(wildcard $(addprefix $(smoke-tests.path), /cases/*))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))

smoke-tests.run: smoke-tests.images
	cd $(smoke-tests.path) && go test -v -tags=smoke

tests smoke-tests .PHONY: smoke-tests.run

template-tests.run: template-tests.images
	cd $(smoke-tests.path) && go test -v -tags=template

tests template-tests .PHONY: template-tests.run
