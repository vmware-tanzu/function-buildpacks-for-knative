RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .
echo-ce.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

echo-ce.image := kn-fn-test/echo-ce
echo-ce.image_paths := $(shell find $(echo-ce.path) -mindepth 1 -maxdepth 1 -type d)
$(echo-ce.image_paths): $(PACK) $(builder.image.out)
	cd $@ && $(PACK) build $(echo-ce.image):$(notdir $@) --builder $(shell cat $(builder.image.out)) --pull-policy if-not-present --clear-cache

echo-ce.clean := $(addsuffix .clean,$(echo-ce.image_paths))
$(echo-ce.clean):
	-docker rmi -f $(echo-ce.image):$(basename $(notdir $@))

.PHONY: smoke-tests.images
smoke-tests.images .PHONY: $(echo-ce.image_paths)
clean .PHONY: $(echo-ce.clean)