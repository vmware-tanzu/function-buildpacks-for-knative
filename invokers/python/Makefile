RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
python-invoker.path := $(path)

python-invoker.sources := $(shell find $(python-invoker.path) -type f -iname '*.py')

python-invoker.build := $(OUT_DIR)/python-invoker.build
$(python-invoker.build): $(python-invoker.sources)
	@mkdir -p $(@D)
	cd $(python-invoker.path) && python setup.py sdist
	touch $@

.PHONY: python-invoker.publish
python-invoker.publish: $(python-invoker.build) # TODO: We should upload to GCP bucket
	@echo TODO: publish Python invoker

.PHONY: python-invoker.tests
python-invoker.tests: $(python-invoker.sources)
	cd $(python-invoker.path) && python -m pytest

.PHONY: python-invoker.clean
python-invoker.clean:
	rm -rf $(python-invoker.path)/.eggs
	rm -rf $(python-invoker.path)/.pytest_cache
	rm -rf $(python-invoker.path)/dist
	rm -rf $(python-invoker.path)/*.egg-info
	rm -rf $(python-invoker.path)/pyfunc/__pycache__

all invokers: $(python-invoker.build)
publish-invokers: python-invoker.publish
tests invoker-tests: python-invoker.tests
clean: python-invoker.clean
