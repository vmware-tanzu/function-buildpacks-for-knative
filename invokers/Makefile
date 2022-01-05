RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
invoker.path := $(abspath $(path))

# Include test case images
invoker.languages := \
	python \
	java \

SUBDIRS := $(addprefix $(invoker.path)/, $(invoker.languages))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))
