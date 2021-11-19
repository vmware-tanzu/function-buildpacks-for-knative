RULES.MK ?= ../../rules.mk
include $(RULES.MK)

PYTHON.MK ?= $(ROOT_DIR)/python.mk
include $(PYTHON.MK)

path ?= .
python-invoker.path := $(path)

python-invoker.sources := $(shell find $(python-invoker.path) -type f -iname '*.py')
python-invoker.version := $(shell cat $(python-invoker.path)/VERSION)

python-invoker.bucket := gs://cnr-functions-release/invokers/python

python-invoker.out_dir := $(OUT_DIR)/invokers/python
python-invoker.build := $(python-invoker.out_dir)/pyfunc-invoker-$(python-invoker.version).tar.gz
$(python-invoker.build): $(TOX) $(python-invoker.sources) $(python-invoker.path)/VERSION $(python-invoker.path)/tox.ini $(python-invoker.path)/Manifest.in
	@mkdir -p $(@D)
	cd $(python-invoker.path) && $(TOX) --sdistonly

python-invoker.sha256 := $(python-invoker.out_dir)/pyfunc-invoker-$(python-invoker.version).sha256
$(python-invoker.sha256): $(python-invoker.build)
	@mkdir -p $(@D)
	cd $(@D) && sha256sum $(notdir $(python-invoker.build)) > $@

python-invoker.artifacts := $(python-invoker.build) $(python-invoker.sha256)
python-invoker.publish := $(OUT_DIR)/python-invoker.publish.$(python-invoker.version)
$(python-invoker.publish): $(GSUTIL) $(python-invoker.artifacts)
	$(GSUTIL) cp $(python-invoker.artifacts) $(python-invoker.bucket)
	$(GSUTIL) retention temp set $(addprefix $(python-invoker.bucket)/,$(notdir $(python-invoker.artifacts)))
	$(GSUTIL) acl ch -u AllUsers:R $(addprefix $(python-invoker.bucket)/,$(notdir $(python-invoker.artifacts)))
	touch $@

.PHONY: python-invoker.tests
python-invoker.tests: $(TOX) $(python-invoker.sources)
	cd $(python-invoker.path) && $(TOX) tests

.PHONY: python-invoker.clean
python-invoker.clean:
	rm -rf $(python-invoker.path)/.eggs
	rm -rf $(python-invoker.path)/.pytest_cache
	rm -rf $(python-invoker.path)/dist
	rm -rf $(python-invoker.path)/*.egg-info
	rm -rf $(python-invoker.path)/pyfunc/__pycache__
	rm -rf $(python-invoker.path)/.tox

.PHONY: publish-invokers.python
publish-invoker.python: $(python-invoker.publish)

.PHONY: invokers.python
invokers.python: $(python-invoker.publish)

all invokers: $(python-invoker.artifacts)
publish-invokers: publish-invokers.python
tests invoker-tests: python-invoker.tests
clean: python-invoker.clean
