RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .

echo-ce.path := $(path)

echo-ce.image := kn-fn-test/echo-ce
echo-ce.image_paths := $(shell find $(echo-ce.path) -mindepth 1 -maxdepth 1 -type d)
$(echo-ce.image_paths): $(PACK)
	-docker rmi -f $(echo-ce.image):$(notdir $@)
	cd $@ && $(PACK) build $(echo-ce.image):$(notdir $@) --builder $(BUILDER_IMAGE)

echo-ce.clean := $(addsuffix .clean,$(echo-ce.image_paths))
$(echo-ce.clean):
	-docker rmi -f $(echo-ce.image):$(basename $(notdir $@))

smoke-test.images .PHONY: $(echo-ce.image_paths)
clean .PHONY: $(echo-ce.clean)
