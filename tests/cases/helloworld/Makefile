RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .
helloworld.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

helloworld.image := kn-fn-test/helloworld
helloworld.image_paths := $(shell find $(helloworld.path) -mindepth 1 -maxdepth 1 -type d)
$(helloworld.image_paths): $(PACK) $(builder.build)
	# -docker rmi -f $(helloworld.image):$(notdir $@)
	cd $@ && $(PACK) build $(helloworld.image):$(notdir $@) --builder $(shell cat $(builder.build))

helloworld.clean := $(addsuffix .clean,$(helloworld.image_paths))
$(helloworld.clean):
	-docker rmi -f $(helloworld.image):$(basename $(notdir $@))

.PHONY: smoke-tests.images
smoke-tests.images .PHONY: $(helloworld.image_paths)
clean .PHONY: $(helloworld.clean)
