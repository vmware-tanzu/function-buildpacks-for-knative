PACK ?= pack
BUILDER_REGISTRY := us.gcr.io/daisy-284300/kn-fn/builder

.PHONY: build
build: builder.toml
	$(PACK) builder create -c builder.toml $(BUILDER_REGISTRY):$(VERSION)

publish: build
	docker push $(BUILDER_REGISTRY):$(VERSION)
