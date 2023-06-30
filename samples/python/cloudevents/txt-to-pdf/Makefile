.PHONY: build publish deploy destroy clean

FUNCTION_IMAGE ?= us.gcr.io/daisy-284300/functions/demo:txttopdf

FUNCTION := out/image
$(FUNCTION): app/*
	@mkdir -p $(@D)
	pack build $(FUNCTION_IMAGE) \
		--path ./app/ \
		--builder paketobuildpacks/builder:0.3.50-base \
		--post-buildpack ghcr.io/vmware-tanzu/function-buildpacks-for-knative/python-buildpack:1.1.2 \
		--env BP_FUNCTION=main.main
	printf $(FUNCTION_IMAGE) > $@

build: $(FUNCTION)

publish: $(FUNCTION)
	docker push $(shell cat $(FUNCTION))

.INTERMEDIATE: $(CONFIG)
CREDS := ./creds.yaml
CONFIG := out/config.yaml
$(CONFIG): $(FUNCTION) $(CREDS)
	@mkdir -p $(@D)
	ytt -f config --data-values-file $(CREDS) -v function_image="$(shell cat $(FUNCTION))" > $@

config: $(CONFIG)

deploy: $(CONFIG) publish
	kapp deploy -a demo -f $(CONFIG)

destroy:
	kapp delete -a demo

clean:
	docker rmi $(FUNCTION_IMAGE) || true
	rm -rf ./out/
