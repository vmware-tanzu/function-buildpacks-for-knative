RULES.MK ?= ../rules.mk
include $(RULES.MK)

ifndef builder.path

path ?= .
builder.path := $(abspath $(path))

$(eval $(call INCLUDE_FILE, $(ROOT_DIR)/buildpacks))

$(eval $(call VERSION_template,builder,$(builder.path)))
builder.version := $(builder.version.release)
builder.build_dir := $(build_dir)/builder/$(builder.version)
builder.out_dir := $(out_dir)/builder/$(builder.version)

####################
# Buildpack Metadata
####################

# If modifying this content make sure you're using SPACES and not TABS for indentation
define builder.buildpack.values.content
#@data/values
---
python_function_buildpack:
  url: $(buildpacks.python.images.offline.name)

java_function_buildpack:
  url: $(buildpacks.java.images.offline.name)
endef

builder.buildpack.values := $(builder.build_dir)/values.yaml

$(builder.buildpack.values): $(ROOT_DIR)/buildpacks/python/VERSION $(ROOT_DIR)/buildpacks/java/VERSION
	@mkdir -p $(@D)
	printf '$(subst $(newline),\n,$(builder.buildpack.values.content))' > $@

##############
# builder.toml
##############

builder.config := $(builder.build_dir)/builder.toml
builder.config_intermediate := $(builder.build_dir)/builder.yaml

builder.config.sources := \
	$(builder.path)/ytt \
	$(shell find '$(builder.path)/ytt' -type f) \

.INTERMEDIATE: $(builder.config_intermediate)
$(builder.config_intermediate): $(YTT) $(builder.config.sources) $(builder.buildpack.values)
	@mkdir -p $(@D)
	$(YTT) -f $(builder.path)/ytt -f $(builder.buildpack.values) --ignore-unknown-comments > $@

$(builder.config): $(YJ) $(builder.config_intermediate)
	@mkdir -p $(@D)
	$(YJ) -yt < $(builder.config_intermediate) > $@

#######
# Build
#######

builder.image.name := $(registry)/functions-builder:$(builder.version)
builder.image.out := $(builder.out_dir)/builder.image

$(builder.image.out): $(PACK) $(builder.config)
	@mkdir -p $(@D)
	$(PACK) builder create -c $(builder.config) $(builder.image.name)
	printf $(builder.image.name) > $@

#########
# Targets
#########

builder.image: $(builder.image.out) $(builder.image.out).sha256

.PHONY: builder.publish
builder.publish: builder.image
	docker push $(builder.image.name)

builder.clean:
	-docker rmi $(builder.image.name)
	rm -rf $(builder.build_dir)
	rm -rf $(builder.out_dir)

all builder .PHONY: builder.image
clean .PHONY: builder.clean

endif
