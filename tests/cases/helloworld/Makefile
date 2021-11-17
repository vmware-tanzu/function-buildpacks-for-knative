RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .

helloworld.image := kn-fn-test/helloworld
helloworld.image_paths := $(shell find $(path) -mindepth 1 -maxdepth 1 -type d)
$(helloworld.image_paths): $(PACK)
	-docker rmi -f $(helloworld.image):$(notdir $@)
	cd $@ && $(PACK) build $(helloworld.image):$(notdir $@) --builder $(BUILDER_IMAGE)

helloworld.clean := $(addsuffix .clean,$(helloworld.image_paths))
$(helloworld.clean):
	-docker rmi -f $(helloworld.image):$(basename $(notdir $@))

smoke-test.images .PHONY: $(helloworld.image_paths)
clean .PHONY: $(helloworld.clean)
