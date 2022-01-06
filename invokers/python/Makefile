RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
invoker.python.path := $(abspath $(path))

include $(invoker.python.path)/../rules.mk

invoker.python.sources := $(shell find $(invoker.python.path) -type f -iname '*.py') \
	$(invoker.python.path)/pyproject.toml \
	$(invoker.python.path)/MANIFEST.in \
	$(invoker.python.path)/VERSION \
	$(invoker.python.path)/LICENSE \
	$(shell find $(invoker.python.path) -mindepth 1 -type d) \

$(eval $(call VERSION_template,invoker.python,$(invoker.python.path)))

# The tools should always build based on the version in VERSION file.
invoker.python.version := $(invoker.python.version.release)

invoker.python.build_dir := $(invoker.build_dir)/python
invoker.python.out_dir := $(invoker.out_dir)/python

invoker.python.artifact.release = $(invoker.python.out_dir)/pyfunc-invoker-$(invoker.python.version).tar.gz

# We want to use the branch here, it's okay as this isn't something we're pushing
invoker.python.builder.tag := $(invoker.python.version.branch)
invoker.python.builder.image_name := python-invoker:$(invoker.python.builder.tag)
invoker.python.builder.image := $(invoker.python.build_dir)/image.$(invoker.python.builder.tag)

$(invoker.python.builder.image): $(invoker.python.path)/Dockerfile $(invoker.python.sources)
	@mkdir -p $(@D)
	DOCKER_BUILDKIT=1 docker build -t $(invoker.python.builder.image_name) -f $< $(invoker.python.path)
	printf $(invoker.python.builder.image_name) > $@

$(invoker.python.artifact.release): $(invoker.python.builder.image)
	@mkdir -p $(@D)
	docker cp `docker create $(invoker.python.builder.image_name)`:/out/$(notdir $(invoker.python.artifact.release)) $@
	touch $@

invoker.python.bucket := $(invoker.bucket)/python

invoker.python.artifacts := $(invoker.python.artifact.release) $(invoker.python.artifact.release).sha256
invoker.python.publish := $(invoker.python.build_dir)/publish.$(invoker.python.version)
$(invoker.python.publish): $(GSUTIL) $(invoker.python.artifacts)
	$(GSUTIL) cp $(invoker.python.artifacts) $(invoker.python.bucket)
	# $(GSUTIL) retention temp set $(addprefix $(invoker.python.bucket)/,$(notdir $(invoker.python.artifacts)))
	$(GSUTIL) acl ch -u AllUsers:R $(addprefix $(invoker.python.bucket)/,$(notdir $(invoker.python.artifacts)))
	touch $@

.PHONY: invoker.python.tests
invoker.python.tests: $(invoker.python.builder.image) $(invoker.python.sources)
	docker run $(shell cat $<)

.PHONY: invoker.python.clean
invoker.python.clean:
	# -docker rmi -f $(invoker.python.builder.image_name)
	$(RM) -r $(invoker.python.build_dir)
	$(RM) -r $(invoker.python.out_dir)

.PHONY: invoker.python invoker.python.publish
invokers invoker.python: $(invoker.python.artifacts)
print-sha: $(invoker.python.artifact.release).print_sha
publish-invokers invoker.python.publish: $(invoker.python.publish)
tests invokers.tests: invoker.python.tests
clean: invoker.python.clean
