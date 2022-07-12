RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
invokers.java.path := $(abspath $(path))

include $(invokers.java.path)/../rules.mk

invokers.java.sources := $(shell find $(invokers.java.path) -type f -iname '*.java') \
	$(invokers.java.path)/pom.xml \
	$(invokers.java.path)/VERSION \

$(eval $(call VERSION_template,invokers.java,$(invokers.java.path)))

# The tools should always build based on the version in VERSION file.
invokers.java.version := $(invokers.java.version.release)

invokers.java.build_dir := $(invokers.build_dir)/java
invokers.java.out_dir := $(invokers.out_dir)/java

invokers.java.artifact.release = $(invokers.java.out_dir)/java-function-invoker-$(invokers.java.version).jar

# We want to use the branch here, it's okay as this isn't something we're pushing
invokers.java.builder.tag := $(invokers.java.version.branch)
invokers.java.builder.image_name := java-invoker:$(invokers.java.builder.tag)
invokers.java.builder.image := $(invokers.java.build_dir)/image.$(invokers.java.builder.tag)

$(invokers.java.builder.image): $(invokers.java.path)/Dockerfile $(invokers.java.sources)
	@mkdir -p $(@D)
	DOCKER_BUILDKIT=1 docker build -t $(invokers.java.builder.image_name) -f $< $(invokers.java.path)
	printf $(invokers.java.builder.image_name) > $@

$(invokers.java.artifact.release): $(invokers.java.builder.image)
	@mkdir -p $(@D)
	docker cp `docker create $(invokers.java.builder.image_name)`:/out/$(notdir $(invokers.java.artifact.release)) $@
	touch $@

invokers.java.bucket := $(invokers.bucket)/java

invokers.java.artifacts := $(invokers.java.artifact.release) $(invokers.java.artifact.release).sha256
invokers.java.publish := $(invokers.java.build_dir)/publish.$(invokers.java.version)
$(invokers.java.publish): $(GSUTIL) $(invokers.java.artifacts)
	$(GSUTIL) cp $(invokers.java.artifacts) $(invokers.java.bucket)
	$(GSUTIL) retention temp set $(addprefix $(invokers.java.bucket)/,$(notdir $(invokers.java.artifacts)))
	$(GSUTIL) acl ch -u AllUsers:R $(addprefix $(invokers.java.bucket)/,$(notdir $(invokers.java.artifacts)))
	touch $@

.PHONY: invokers.java.tests
invokers.java.tests: $(invokers.java.builder.image) $(invokers.java.sources)
	docker run $(shell cat $<)

.PHONY: invokers.java.clean
invokers.java.clean:
	# -docker rmi -f $(invokers.java.builder.image_name)
	$(RM) -r $(invokers.java.build_dir)
	$(RM) -r $(invokers.java.out_dir)

.PHONY: invokers.java invokers.java.publish
invokers invokers.java: $(invokers.java.artifacts)
invokers.publish invokers.java.publish: $(invokers.java.publish)
print-sha invokers.java.print-sha: $(invokers.java.artifact.release).print-sha
tests invokers.tests: invokers.java.tests
clean: invokers.java.clean
