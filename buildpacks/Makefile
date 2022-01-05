RULES.MK ?= ../rules.mk
include $(RULES.MK)

ifndef buildpack.path

path ?= .
buildpack.path := $(abspath $(path))

# Include test case images
buildpack.languages = \
	python \
	java \

SUBDIRS := $(addprefix $(buildpack.path)/, $(buildpack.languages))
$(foreach dir,$(SUBDIRS),$(eval $(call INCLUDE_FILE, $(dir))))
endif
