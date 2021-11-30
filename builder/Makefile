RULES.MK ?= ../rules.mk
include $(RULES.MK)

path ?= .
builder.path := $(path)

builder.build := $(abspath $(OUT_DIR)/builder.build)
$(builder.build): $(PACK) $(builder.path)/builder.toml $(foreach lang,$(buildpack.languages),$($(lang)-buildpack.deps))
	@mkdir -p $(@D)
	$(PACK) builder create -c $(builder.path)/builder.toml $(BUILDER_IMAGE)
	touch $@

builder.publish: $(builder.build)
	docker push $(BUILDER_IMAGE)

builder.clean:
	-docker rmi $(BUILDER_IMAGE)
	rm -f $(builder.build)

all builder: $(builder.build)
publish publish-builder .PHONY: builder.publish
clean .PHONY: builder.clean
