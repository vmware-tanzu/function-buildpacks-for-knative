# After including this file, add the following line to have the version printed
# when invoking the `version` target:
# $(eval $(call VERSION_template,target_name,path_to_VERSION_file))
# To override any of the values, redefine it after the above line.
define VERSION_template
$(1).version.release := $(shell cat $(2)/VERSION)
$(1).version.dev := dev-$(build_time)
$(1).version.commit := dev-$(git.commit)
$(1).version.branch := dev-$(subst /,_,$(git.branch))
$(1).version.latest := latest
$(1).version = $$($(1).version.$(build))
$(1).version:
	@echo $(1): $$($(1).version)

version .PHONY: $(1).version
endef
