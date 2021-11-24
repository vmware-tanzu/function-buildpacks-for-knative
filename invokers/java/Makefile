RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
java-invoker.path := $(path)

java-invoker.sources := $(shell find $(java-invoker.path) -type f -iname '*.java') \
	$(java-invoker.path)/pom.xml \
	$(java-invoker.path)/VERSION \

java-invoker.version := $(shell cat $(java-invoker.path)/VERSION)

java-invoker.bucket := gs://cnr-functions-release/invokers/java

java-invoker.build_dir := $(BUILD_DIR)/java-invoker
java-invoker.build_image_name := java-invoker:$(GIT_COMMIT)
java-invoker.build_image := $(java-invoker.build_dir)/image.$(GIT_COMMIT)

java-invoker.out_dir := $(OUT_DIR)/invokers/java
java-invoker.out_file := $(java-invoker.out_dir)/java-function-invoker-$(java-invoker.version).jar

$(java-invoker.build_image): $(java-invoker.path)/Dockerfile $(java-invoker.sources)
	@mkdir -p $(@D)
	DOCKER_BUILDKIT=1 docker build -t $(java-invoker.build_image_name) -f $< $(java-invoker.path)
	touch $@

$(java-invoker.out_file): $(java-invoker.build_image)
	@mkdir -p $(@D)
	docker cp `docker create $(java-invoker.build_image_name)`:/out/$(notdir $(java-invoker.out_file)) $@
	touch $@

java-invoker.artifacts := $(java-invoker.out_file) $(java-invoker.out_file).sha256
java-invoker.publish := $(java-invoker.build_dir)/publish.$(java-invoker.version)
$(java-invoker.publish): $(GSUTIL) $(java-invoker.artifacts)
	$(GSUTIL) cp $(java-invoker.artifacts) $(java-invoker.bucket)
	$(GSUTIL) retention temp set $(addprefix $(java-invoker.bucket)/,$(notdir $(java-invoker.artifacts)))
	$(GSUTIL) acl ch -u AllUsers:R $(addprefix $(java-invoker.bucket)/,$(notdir $(java-invoker.artifacts)))
	touch $@

.PHONY: java-invoker.print_sha
java-invoker.print_sha: $(java-invoker.out_file).print_sha

.PHONY: java-invoker.tests
java-invoker.tests: $(java-invoker.build_image) $(java-invoker.sources)
	docker run $(java-invoker.build_image_name)
	@echo Java invoker missing test implementation.

.PHONY: java-invoker.clean
java-invoker.clean:
	-docker rmi -f $(java-invoker.build_image_name)

.PHONY: invokers.java publish-invokers.java
publish-invoker.java: $(java-invoker.publish)

all invokers invokers.java: $(java-invoker.artifacts)
publish-invokers: publish-invoker.java
print-sha: java-invoker.print_sha
tests invoker-tests: java-invoker.tests
clean: java-invoker.clean
