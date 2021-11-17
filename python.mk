ifndef PYTHON_MK # Prevent repeated "-include".
PYTHON_MK := $(lastword $(MAKEFILE_LIST))
ROOT_DIR := $(dir $(PYTHON_MK))
include $(ROOT_DIR)/rules.mk

python-venv.dir := $(OUT_DIR)/python-venv
python-venv.install := $(OUT_DIR)/python-venv.install

PYTHON := $(python-venv.dir)/bin/python
TOX := $(python-venv.dir)/bin/tox

# This will create a virtual python environment
# Within this environment is the python binary we should use.
# Thus we can't use $(PYTHON) here
$(python-venv.dir):
	python -m venv $@

$(python-venv.install).%: $(python-venv.dir)
	$(PYTHON) -m pip install $*
	touch $@

$(PYTHON): $(python-venv.dir)

# Python Tools
$(TOX): $(python-venv.install).tox

python-venv.clean:
	rm -rf $(python-venv.dir)

clean: python-venv.clean
endif
