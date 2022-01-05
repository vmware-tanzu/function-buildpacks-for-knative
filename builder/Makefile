RULES.MK ?= ../rules.mk
include $(RULES.MK)

ifndef builder.path

path ?= .
builder.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/buildpacks))

$(eval $(call VERSION_template,builder,$(builder.path)))
builder.image := $(registry)/builder:$(builder.version)
builder.image_full := $(builder.image)

builder.build := $(abspath $(out_dir)/builder.build)
$(builder.build): $(PACK) $(builder.path)/builder.toml $(foreach lang,$(buildpack.languages),$($(lang)-buildpack.deps))
	@mkdir -p $(@D)
	$(PACK) builder create -c $(builder.path)/builder.toml $(builder.image)
	printf $(builder.image_full) > $@

builder.publish: $(builder.build)
	docker push $(builder.image)

builder.clean:
	-docker rmi $(builder.image)
	rm -f $(builder.build)

all builder: $(builder.build)
publish-builder .PHONY: builder.publish
clean .PHONY: builder.clean

endif
