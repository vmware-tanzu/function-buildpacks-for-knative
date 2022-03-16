RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
smoke-tests.path := $(abspath $(path))

# Include test case images
SUBDIRS := $(wildcard $(addprefix $(smoke-tests.path), /cases/*))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))

smoke-tests.run: smoke-tests.images
	cd $(smoke-tests.path) && go test -run TestSmoke

tests smoke-tests .PHONY: smoke-tests.run

template-tests.run: smoke-tests.images
	cd $(smoke-tests.path) && go test -run TestTemplates

tests template-tests .PHONY: template-tests.run
