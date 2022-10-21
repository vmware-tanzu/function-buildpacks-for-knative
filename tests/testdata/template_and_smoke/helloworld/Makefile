RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .
helloworld.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

helloworld.env.java := com.example.demo.HelloWorld
helloworld.env.python := handler.handler

helloworld.image := kn-fn-test/helloworld
helloworld.image_paths := $(filter-out %/Makefile, $(wildcard $(helloworld.path)/*))

$(helloworld.image_paths): $(PACK) $(builder.image.out)
	cd $@ && $(PACK) build $(helloworld.image):$(notdir $@) --builder $(shell cat $(builder.image.out)) --env BP_FUNCTION=$(helloworld.env.$(notdir $@)) --pull-policy if-not-present --clear-cache

helloworld.clean := $(addsuffix .clean,$(helloworld.image_paths))
$(helloworld.clean):
	-docker rmi -f $(helloworld.image):$(basename $(notdir $@))

.PHONY: smoke-tests.images
smoke-tests.images .PHONY: $(helloworld.image_paths)
clean .PHONY: $(helloworld.clean)
