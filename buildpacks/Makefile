RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .

# Include test case images
buildpack.languages = \
	python \
	java \

SUBDIRS := $(addprefix $(path)/, $(buildpack.languages))
$(foreach dir,$(SUBDIRS),$(eval $(INCLUDE_FILE)))
