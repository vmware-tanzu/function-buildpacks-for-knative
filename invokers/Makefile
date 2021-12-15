RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .

# Include test case images
invoker.languages := \
	python \
	java \

SUBDIRS := $(addprefix $(path)/, $(invoker.languages))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))
