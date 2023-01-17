# -*- makefile -*-
## -----------------------------------------------------------------------
## -----------------------------------------------------------------------

repos += voltha-system-tests
repos += voltha-helm-charts
repos += pod-configs
repos += voltha-docs



repo-deps = $(addprefix sandbox/,$(repos))

get-arg-one = $(lastword $(subst /,$(space),$(1)))
checkout-repo = ssh://gerrit.opencord.org:29418/$(call get-arg-one,$(1)).git

$(repo-deps):
	@mkdir -p $(dir $@)
	@cd sandbox && $(GIT) clone $(call checkout-repo,$@)

checkout-repos-all : $(repo-deps)

.PHONY: sandbox
sandbox : $(repo-deps)

# [EOF]
