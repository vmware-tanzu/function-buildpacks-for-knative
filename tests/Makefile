RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
smoke-tests.path := $(path)

# Include test case images
SUBDIRS := $(wildcard $(addprefix $(path), /cases/*))
$(foreach dir,$(SUBDIRS),$(eval $(INCLUDE_FILE)))

smoke-tests.run: smoke-tests.images
	cd $(smoke-tests.path) && go test -v

tests smoke-tests .PHONY: smoke-tests.run
