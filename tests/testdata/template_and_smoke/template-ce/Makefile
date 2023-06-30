RULES.MK ?= ../../../../rules.mk
include $(RULES.MK)

path ?= .
template-ce.path := $(abspath $(path))

template-ce.templates.langs := $(filter-out %/README.md, $(wildcard $(abspath $(ROOT_DIR)/templates)/*))
template-ce.sources := $(foreach lang, $(template-ce.templates.langs), $(addprefix $(lang)/, $(filter cloudevents%, $(notdir $(wildcard $(lang)/*)))))
template-ce.dest := $(addprefix $(template-ce.path)/, $(foreach lang,$(template-ce.templates.langs),$(addprefix $(notdir $(lang))-, $(filter cloudevents%, $(notdir $(wildcard $(lang)/*))))))

template-ce.out := $(addprefix $(out_dir)/tests/template/,$(notdir $(template-ce.dest)))

template-ce.java.bp_function := functions.Handler
template-ce.python.bp_function := func.main

.SECONDEXPANSION:
$(template-ce.path)/%: $$(foreach v,%, $$(addprefix $(ROOT_DIR)/templates/,$$(firstword $$(subst -, ,$$(v)))/$$(subst $$(firstword $$(subst -, ,$$(v)))-,,$$(v))))
	ln -sf $< $@

$(template-ce.out): $(out_dir)/tests/template/%: $(template-ce.path)/% $(PACK)
	cd $< && $(PACK) build $(template-ce.image):$(notdir $@) --builder paketobuildpacks/builder:base --env BP_JVM_VERSION=17 --env BP_FUNCTION=$(template-ce.$(firstword $(subst -, ,$*)).bp_function) --pull-policy if-not-present --clear-cache

template-ce.image := kn-fn-test/template-ce
template-ce.clean := $(addsuffix .clean,$(template-ce.out))
$(template-ce.clean):
	-docker rmi -f $(template-ce.image):$(basename $(notdir $@))

.PHONY: template-tests.images
template-tests.images .PHONY: $(template-ce.out)
clean .PHONY: $(template-ce.clean)
