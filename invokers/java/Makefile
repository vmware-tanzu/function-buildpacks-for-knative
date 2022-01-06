RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
invoker.java.path := $(abspath $(path))

include $(invoker.java.path)/../rules.mk

invoker.java.sources := $(shell find $(invoker.java.path) -type f -iname '*.java') \
	$(invoker.java.path)/pom.xml \
	$(invoker.java.path)/VERSION \

$(eval $(call VERSION_template,invoker.java,$(invoker.java.path)))

# The tools should always build based on the version in VERSION file.
invoker.java.version := $(invoker.java.version.release)

invoker.java.build_dir := $(invoker.build_dir)/java
invoker.java.out_dir := $(invoker.out_dir)/java

invoker.java.artifact.release = $(invoker.java.out_dir)/java-function-invoker-$(invoker.java.version).jar

# We want to use the branch here, it's okay as this isn't something we're pushing
invoker.java.builder.tag := $(invoker.java.version.branch)
invoker.java.builder.image_name := java-invoker:$(invoker.java.builder.tag)
invoker.java.builder.image := $(invoker.java.build_dir)/image.$(invoker.java.builder.tag)

$(invoker.java.builder.image): $(invoker.java.path)/Dockerfile $(invoker.java.sources)
	@mkdir -p $(@D)
	DOCKER_BUILDKIT=1 docker build -t $(invoker.java.builder.image_name) -f $< $(invoker.java.path)
	printf $(invoker.java.builder.image_name) > $@

$(invoker.java.artifact.release): $(invoker.java.builder.image)
	@mkdir -p $(@D)
	docker cp `docker create $(invoker.java.builder.image_name)`:/out/$(notdir $(invoker.java.artifact.release)) $@
	touch $@

invoker.java.bucket := $(invoker.bucket)/java

invoker.java.artifacts := $(invoker.java.artifact.release) $(invoker.java.artifact.release).sha256
invoker.java.publish := $(invoker.java.build_dir)/publish.$(invoker.java.version)
$(invoker.java.publish): $(GSUTIL) $(invoker.java.artifacts)
	$(GSUTIL) cp $(invoker.java.artifacts) $(invoker.java.bucket)
	# $(GSUTIL) retention temp set $(addprefix $(invoker.java.bucket)/,$(notdir $(invoker.java.artifacts)))
	$(GSUTIL) acl ch -u AllUsers:R $(addprefix $(invoker.java.bucket)/,$(notdir $(invoker.java.artifacts)))
	touch $@

.PHONY: invoker.java.tests
invoker.java.tests: $(invoker.java.builder.image) $(invoker.java.sources)
	docker run $(shell cat $<)

.PHONY: invoker.java.clean
invoker.java.clean:
	# -docker rmi -f $(invoker.java.builder.image_name)
	$(RM) -r $(invoker.java.build_dir)
	$(RM) -r $(invoker.java.out_dir)

.PHONY: invoker.java invoker.java.publish
invokers invoker.java: $(invoker.java.artifacts)
print-sha: $(invoker.java.artifact.release).print_sha
publish-invokers invoker.java.publish: $(invoker.java.publish)
tests invokers.tests: invoker.java.tests
clean: invoker.java.clean
