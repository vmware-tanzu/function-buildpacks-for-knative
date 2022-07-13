RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .
template-http.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

template-http.image := kn-fn-test/template-http
template-http.image_paths := $(shell find $(template-http.path) -mindepth 1 -maxdepth 1 -type l)
$(template-http.image_paths): $(PACK) $(builder.image.out)
	cd $@ && $(PACK) build $(template-http.image):$(notdir $@) --builder $(shell cat $(builder.image.out)) --pull-policy if-not-present --clear-cache

template-http.clean := $(addsuffix .clean,$(template-http.image_paths))
$(template-http.clean):
	-docker rmi -f $(template-http.image):$(basename $(notdir $@))

.PHONY: template-tests.images
template-tests.images .PHONY: $(template-http.image_paths)
clean .PHONY: $(template-http.clean)
