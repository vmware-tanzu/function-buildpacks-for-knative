RULES.MK ?= ../../../rules.mk
include $(RULES.MK)

path ?= .
python-invoker.path := $(path)

python-invoker.sources := $(shell find $(python-invoker.path) -type f -iname '*.py')

python-invoker.build := $(OUT_DIR)/python-invoker.build
$(python-invoker.build): $(python-invoker.sources)
	@mkdir -p $(@D)
	python setup.py sdist
	touch $@

python-invoker.publish:; # TODO: We should upload to GCP bucket

python-invoker.tests: $(python-invoker.sources)
	cd $(python-invoker.path) && python -m pytest

python-invoker.clean:
	rm -rf $(python-invoker.path)/.eggs
	rm -rf $(python-invoker.path)/.pytest_cache
	rm -rf $(python-invoker.path)/dist
	rm -rf $(python-invoker.path)/*.egg-info
	rm -rf $(python-invoker.path)/pyfunc/__pycache__

# This is part of the default rules
all .PHONY: python-invoker.publish
publish-invokers .PHONY: python-invoker.publish
tests .PHONY: python-invoker.tests
clean .PHONY: python-invoker.clean
