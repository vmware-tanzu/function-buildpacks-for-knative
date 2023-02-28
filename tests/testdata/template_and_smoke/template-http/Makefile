RULES.MK ?= ../../../../rules.mk
include $(RULES.MK)

path ?= .
template-http.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/builder))

template-http.templates.langs := $(filter-out %/README.md, $(wildcard $(abspath $(ROOT_DIR)/templates)/*))
template-http.sources := $(foreach lang, $(template-http.templates.langs), $(addprefix $(lang)/, $(filter http%, $(notdir $(wildcard $(lang)/*)))))
template-http.dest := $(addprefix $(template-http.path)/, $(foreach lang,$(template-http.templates.langs),$(addprefix $(notdir $(lang))-, $(filter http%, $(notdir $(wildcard $(lang)/*))))))

template-http.out := $(addprefix $(out_dir)/tests/template/,$(notdir $(template-http.dest)))

template-http.java.bp_function := functions.Handler
template-http.python.bp_function := func.main

.SECONDEXPANSION:
$(template-http.path)/%: $$(foreach v,%, $$(addprefix $(ROOT_DIR)/templates/,$$(firstword $$(subst -, ,$$(v)))/$$(subst $$(firstword $$(subst -, ,$$(v)))-,,$$(v))))
	ln -sf $< $@

$(template-http.out): $(out_dir)/tests/template/%: $(template-http.path)/% $(PACK) $(builder.image.out)
	cd $< && $(PACK) build $(template-http.image):$(notdir $@) --builder $(shell cat $(builder.image.out)) --env BP_FUNCTION=$(template-http.$(firstword $(subst -, ,$*)).bp_function) --env BP_JVM_VERSION=17 --pull-policy if-not-present --clear-cache

template-http.image := kn-fn-test/template-http
template-http.clean := $(addsuffix .clean,$(template-http.out))
$(template-http.clean):
	-docker rmi -f $(template-http.image):$(basename $(notdir $@))

.PHONY: template-tests.images
template-tests.images .PHONY: $(template-http.out)
clean .PHONY: $(template-http.clean)
