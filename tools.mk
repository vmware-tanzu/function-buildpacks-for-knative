ifndef TOOLS_MK # Prevent repeated "-include".
TOOLS_MK := $(lastword $(MAKEFILE_LIST))
TOOLS_INCLUDE_DIR := $(dir $(TOOLS_MK))

# Define the tools here
tools.path := $(abspath $(build_dir)/tools)
tools.bin.path := $(abspath $(tools.path)/bin)

PACK := $(tools.bin.path)/pack
CREATE-PACKAGE := $(tools.bin.path)/create-package
GSUTIL := $(tools.path)/gsutil/gsutil

$(PACK).darwin:
	@mkdir -p $(@D)
	curl -sL https://github.com/buildpacks/pack/releases/download/v0.21.1/pack-v0.21.1-macos.tgz | tar -xz -C $(@D)
	touch $@

$(PACK).linux:
	@mkdir -p $(@D)
	curl -sL https://github.com/buildpacks/pack/releases/download/v0.21.1/pack-v0.21.1-linux.tgz | tar -xz -C $(@D)
	touch $@

$(PACK): $(PACK).$(os.name)
	chmod +x $@
	touch $@

$(CREATE-PACKAGE):
	@mkdir -p $(@D)
	GOBIN=$(@D) go install github.com/paketo-buildpacks/libpak/cmd/create-package@latest

$(GSUTIL):
	@mkdir -p $(@D)
	curl -sL https://storage.googleapis.com/pub/gsutil.tar.gz | tar -xz -C $(tools.path)

tools.clean:
	$(RM) -rf $(tools.path)

clean .PHONY: tools.clean

endif
