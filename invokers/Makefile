RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
invokers.path := $(abspath $(path))

include $(invokers.path)/rules.mk

# Include test case images
invokers.languages := \
	python \
	java \

SUBDIRS := $(addprefix $(invokers.path)/, $(invokers.languages))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))
