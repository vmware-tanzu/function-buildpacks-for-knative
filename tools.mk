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

YTT := $(tools.bin.path)/ytt
YTT.version := v0.41.1

$(YTT).darwin:
	@mkdir -p $(@D)
	curl -sL https://github.com/vmware-tanzu/carvel-ytt/releases/download/$(YTT.version)/ytt-darwin-amd64 -o $@
	chmod +x $@
	touch $@

$(YTT).linux:
	@mkdir -p $(@D)
	curl -sL https://github.com/vmware-tanzu/carvel-ytt/releases/download/$(YTT.version)/ytt-linux-amd64 -o $@
	chmod +x $@
	touch $@

$(YTT): $(YTT).$(os.name)
	ln -sf $< $@

YJ := $(tools.bin.path)/yj
YJ.version := v5.1.0
$(YJ).darwin:
	@mkdir -p $(@D)
	curl -sL https://github.com/sclevine/yj/releases/download/$(YJ.version)/yj-macos-amd64 -o $@
	chmod +x $@
	touch $@

$(YJ).linux:
	@mkdir -p $(@D)
	curl -sL https://github.com/sclevine/yj/releases/download/$(YJ.version)/yj-linux-amd64 -o $@
	chmod +x $@
	touch $@

$(YJ): $(YJ).$(os.name)
	ln -sf $< $@

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
