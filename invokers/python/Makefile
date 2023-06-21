RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
invokers.python.path := $(abspath $(path))

include $(invokers.python.path)/../rules.mk

invokers.python.sources := $(shell find $(invokers.python.path) -type f -iname '*.py') \
	$(invokers.python.path)/pyproject.toml \
	$(invokers.python.path)/MANIFEST.in \
	$(invokers.python.path)/VERSION \
	$(shell find $(invokers.python.path) -mindepth 1 -type d) \

$(eval $(call VERSION_template,invokers.python,$(invokers.python.path)))

# The tools should always build based on the version in VERSION file.
invokers.python.version := $(invokers.python.version.release)

invokers.python.build_dir := $(invokers.build_dir)/python
invokers.python.out_dir := $(invokers.out_dir)/python

invokers.python.artifact.release = $(invokers.python.out_dir)/pyfunc-invoker-$(invokers.python.version).tar.gz
invokers.python.artifact.release-deps = $(invokers.python.out_dir)/pyfunc-invoker-deps-$(invokers.python.version).tar.gz

# We want to use the branch here, it's okay as this isn't something we're pushing
invokers.python.builder.tag := $(invokers.python.version.branch)
invokers.python.builder.image_name := python-invoker:$(invokers.python.builder.tag)
invokers.python.builder.image := $(invokers.python.build_dir)/image.$(invokers.python.builder.tag)

$(invokers.python.builder.image): $(invokers.python.path)/Dockerfile $(invokers.python.sources)
	@mkdir -p $(@D)
	DOCKER_BUILDKIT=1 docker build -t $(invokers.python.builder.image_name) -f $< $(invokers.python.path)
	printf $(invokers.python.builder.image_name) > $@

$(invokers.python.artifact.release): $(invokers.python.builder.image)
	@mkdir -p $(@D)
	docker cp `docker create $(invokers.python.builder.image_name)`:/out/$(notdir $(invokers.python.artifact.release)) $@
	touch $@

$(invokers.python.artifact.release-deps): $(invokers.python.builder.image)
	@mkdir -p $(@D)
	docker cp `docker create $(invokers.python.builder.image_name)`:/out/$(notdir $(invokers.python.artifact.release-deps)) $@
	touch $@

invokers.python.bucket := $(invokers.bucket)/python

invokers.python.artifacts := $(invokers.python.artifact.release) $(invokers.python.artifact.release).sha256 $(invokers.python.artifact.release-deps) $(invokers.python.artifact.release-deps).sha256
invokers.python.publish := $(invokers.python.build_dir)/publish.$(invokers.python.version)
$(invokers.python.publish): $(GSUTIL) $(invokers.python.artifacts)
	$(GSUTIL) cp $(invokers.python.artifacts) $(invokers.python.bucket)
	$(GSUTIL) retention temp set $(addprefix $(invokers.python.bucket)/,$(notdir $(invokers.python.artifacts)))
	$(GSUTIL) acl ch -u AllUsers:R $(addprefix $(invokers.python.bucket)/,$(notdir $(invokers.python.artifacts)))
	touch $@

.PHONY: invokers.python.tests
invokers.python.tests: $(invokers.python.builder.image) $(invokers.python.sources)
	docker run -v $(invokers.python.path)/:/workspace/invoker/ $(shell cat $<)

.PHONY: invokers.python.clean
invokers.python.clean:
	$(RM) -r $(invokers.python.build_dir)
	$(RM) -r $(invokers.python.out_dir)

invokers.python: $(invokers.python.artifacts)
invokers.python.publish: $(invokers.python.publish)
invokers.python.print-sha: $(invokers.python.artifact.release).print-sha

invokers.publish .PHONY: invokers.python.publish
invokers.print-sha .PHONY: invokers.python.print-sha
invokers.tests .PHONY: invokers.python.tests
invokers.clean .PHONY: invokers.python.clean
invokers .PHONY: invokers.python
