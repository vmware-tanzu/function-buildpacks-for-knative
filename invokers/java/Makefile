RULES.MK ?= ../../rules.mk
include $(RULES.MK)

path ?= .
java-invoker.path := $(path)

java-invoker.sources := $(shell find $(java-invoker.path) -type f -iname '*.py')

java-invoker.build := $(OUT_DIR)/java-invoker.build
$(java-invoker.build): $(java-invoker.sources)
	@echo TODO: implement build for Java invoker. It might be just "mvn package"

.PHONY: java-invoker.publish
java-invoker.publish: $(java-invoker.build) # TODO: We should upload to GCP bucket
	@echo TODO: publish Java invoker

.PHONY: java-invoker.tests
java-invoker.tests: $(java-invoker.sources)
	@echo TODO: publish Java invoker tests

.PHONY: java-invoker.clean
java-invoker.clean:

all invokers: $(java-invoker.build)
publish-invokers: java-invoker.publish
tests invoker-tests: java-invoker.tests
clean: java-invoker.clean
