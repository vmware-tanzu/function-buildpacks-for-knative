RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .
template-ce.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

template-ce.image := kn-fn-test/template-ce
template-ce.java_image_paths := $(shell find $(template-ce.path) -mindepth 1 -maxdepth 1 -type l | grep java-)
template-ce.python_image_paths := $(shell find $(template-ce.path) -mindepth 1 -maxdepth 1 -type l | grep python-)
$(template-ce.java_image_paths): $(PACK) $(builder.image.out)
	cd $@ && $(PACK) build $(template-ce.image):$(notdir $@) --builder $(shell cat $(builder.image.out)) --env BP_FUNCTION=functions.Handler --pull-policy if-not-present --clear-cache
$(template-ce.python_image_paths): $(PACK) $(builder.image.out)
	cd $@ && $(PACK) build $(template-ce.image):$(notdir $@) --builder $(shell cat $(builder.image.out)) --env BP_FUNCTION=func.main --pull-policy if-not-present --clear-cache

template-ce.clean := $(addsuffix .clean,$(template-ce.image_paths))
$(template-ce.clean):
	-docker rmi -f $(template-ce.image):$(basename $(notdir $@))

.PHONY: template-tests.images
template-tests.images .PHONY: $(template-ce.java_image_paths) $(template-ce.python_image_paths)
clean .PHONY: $(template-ce.clean)
