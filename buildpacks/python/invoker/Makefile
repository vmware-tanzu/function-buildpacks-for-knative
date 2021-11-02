.PHONY: buildpack
buildpack:
	cd buildpack && pack buildpack package vmware/pyfun-addin:0.1.0

.PHONY: metabuild
metabuild: buildpack
	cd metabuild && pack buildpack package vmware/pyfun:0.1.1 --config ./package.toml

.PHONY: http-func
http-func:
	cd functions/pyfun/http && pack build pytestapp --buildpack vmware/pyfun:0.1.1 --builder paketobuildpacks/builder:base

.PHONY: run
run:
	docker run --rm -p 8080:8080 -e 8080 pytestapp

test:
	pytest
