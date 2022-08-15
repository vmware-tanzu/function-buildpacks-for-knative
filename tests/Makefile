RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
tests.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

# Include test case images
SUBDIRS := $(wildcard $(addprefix $(tests.path), /testdata/template_and_smoke/*))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))

smoke-tests.run: smoke-tests.images
	cd $(tests.path) && go test -v -tags=smoke

tests smoke-tests .PHONY: smoke-tests.run

template-tests.run: template-tests.images
	cd $(tests.path) && go test -v -tags=template

tests template-tests .PHONY: template-tests.run
