RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
smoke-test.path := $(path)

# Include test case images
SUBDIRS := $(wildcard $(addprefix $(path), /cases/*))
$(foreach dir,$(SUBDIRS),$(eval $(INCLUDE_FILE)))

smoke-tests.run: smoke-test.images
	cd $(smoke-test.path) && go test -v

tests smoke-tests .PHONY: smoke-tests.run
