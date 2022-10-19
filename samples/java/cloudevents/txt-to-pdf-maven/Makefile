.PHONY: build publish deploy destroy clean

FUNCTION_IMAGE ?= us.gcr.io/daisy-284300/functions/demo:txttopdfjavamaven

FUNCTION := out/image
$(FUNCTION):
	@mkdir -p $(@D)
	pack build $(FUNCTION_IMAGE) --path ./spring-java-fn/ --builder ghcr.io/vmware-tanzu/function-buildpacks-for-knative/functions-builder:0.1.0 --env BP_FUNCTION=com.example.demo.DemoApplication
	printf $(FUNCTION_IMAGE) > $@

build: $(FUNCTION)

publish: $(FUNCTION)
	docker push $(shell cat $(FUNCTION))

.INTERMEDIATE: $(CONFIG)
CREDS := ./creds.yaml
CONFIG := out/config.yaml
$(CONFIG): $(FUNCTION) $(CREDS) $(wildcard config/*)
	@mkdir -p $(@D)
	ytt -f config --data-values-file $(CREDS) -v function_image="$(shell cat $(FUNCTION))" > $@

config: $(CONFIG)

deploy: $(CONFIG) publish
	kapp deploy -a demo -f $(CONFIG)

destroy:
	kapp delete -a demo

clean:
	docker rmi $(FUNCTION_IMAGE)
	rm -rf ./out/

docker: $(FUNCTION)
	docker run --env-file aws.env --rm -p 8080:8080 -e 8080 --entrypoint function $(FUNCTION_IMAGE)
