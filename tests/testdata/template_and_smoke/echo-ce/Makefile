RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .
echo-ce.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

echo-ce.env.java := com.vmware.functions.Func
echo-ce.env.python := handler.handler

echo-ce.image := kn-fn-test/echo-ce
echo-ce.image_paths := $(filter-out %/Makefile, $(wildcard $(echo-ce.path)/*))

$(echo-ce.image_paths): $(PACK) $(builder.image.out)
	cd $@ && $(PACK) build $(echo-ce.image):$(notdir $@) --builder $(shell cat $(builder.image.out)) --env BP_FUNCTION=$(echo-ce.env.$(notdir $@)) --pull-policy if-not-present --clear-cache

echo-ce.clean := $(addsuffix .clean,$(echo-ce.image_paths))
$(echo-ce.clean):
	-docker rmi -f $(echo-ce.image):$(basename $(notdir $@))

.PHONY: smoke-tests.images
smoke-tests.images .PHONY: $(echo-ce.image_paths)
clean .PHONY: $(echo-ce.clean)
