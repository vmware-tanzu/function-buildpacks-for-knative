IMAGE_NAME := kn-fn-test/helloworld

paths := $(abspath $(shell find ~+ -maxdepth 1 -type d))
$(paths):
	-docker rmi -f $(IMAGE_NAME):$(notdir $@)
	cd $@ && $(PACK) build $(IMAGE_NAME):$(notdir $@) --builder $(BUILDER_IMAGE)

images: $(paths)

.PHONY: images $(paths)
