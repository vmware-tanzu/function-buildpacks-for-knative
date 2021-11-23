RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
python-invoker.path := $(path)

python-invoker.sources := $(shell find $(python-invoker.path) -type f -iname '*.py') \
	$(python-invoker.path)/pyproject.toml \
	$(python-invoker.path)/MANIFEST.in \
	$(python-invoker.path)/VERSION \
	$(python-invoker.path)/LICENSE \

python-invoker.version := $(shell cat $(python-invoker.path)/VERSION)

python-invoker.bucket := gs://cnr-functions-release/invokers/python

python-invoker.build_dir := $(BUILD_DIR)/python-invoker
python-invoker.build_image_name := python-invoker:$(GIT_COMMIT)
python-invoker.build_image := $(python-invoker.build_dir)/image.$(GIT_COMMIT)

python-invoker.out_dir := $(OUT_DIR)/invokers/python
python-invoker.out_file := $(python-invoker.out_dir)/pyfunc-invoker-$(python-invoker.version).tar.gz

$(python-invoker.build_image): $(python-invoker.path)/Dockerfile $(python-invoker.sources)
	@mkdir -p $(@D)
	DOCKER_BUILDKIT=1 docker build -t $(python-invoker.build_image_name) -f $< $(python-invoker.path)
	touch $@

$(python-invoker.out_file): $(python-invoker.build_image)
	@mkdir -p $(@D)
	docker cp `docker create $(python-invoker.build_image_name)`:/out/$(notdir $(python-invoker.out_file)) $@
	touch $@

python-invoker.artifacts := $(python-invoker.out_file) $(python-invoker.out_file).sha256
python-invoker.publish := $(python-invoker.build_dir)/publish.$(python-invoker.version)
$(python-invoker.publish): $(GSUTIL) $(python-invoker.artifacts)
	$(GSUTIL) cp $(python-invoker.artifacts) $(python-invoker.bucket)
	$(GSUTIL) retention temp set $(addprefix $(python-invoker.bucket)/,$(notdir $(python-invoker.artifacts)))
	$(GSUTIL) acl ch -u AllUsers:R $(addprefix $(python-invoker.bucket)/,$(notdir $(python-invoker.artifacts)))
	touch $@

.PHONY: python-invoker.print_sha
python-invoker.print_sha: $(python-invoker.out_file).print_sha

.PHONY: python-invoker.tests
python-invoker.tests: $(python-invoker.build_image) $(python-invoker.sources)
	docker run $(python-invoker.build_image_name)

.PHONY: python-invoker.clean
python-invoker.clean:
	-docker rmi -f $(python-invoker.build_image_name)

.PHONY: invokers.python publish-invokers.python
publish-invoker.python: $(python-invoker.publish)

all invokers invokers.python: $(python-invoker.artifacts)
publish-invokers: publish-invokers.python
print-sha: python-invoker.print_sha
tests invoker-tests: python-invoker.tests
clean: python-invoker.clean
