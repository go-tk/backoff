override SHELL := bash
override .SHELLFLAGS := -o errexit -o nounset -o pipefail -c

override targets := $(or $(MAKECMDGOALS),placeholder)
override build_repo := https://github.com/go-tk/build.git
override build_dir := .build
override build_ttl := 24 hours

.PHONY: $(targets)
.ONESHELL:
$(targets):
	@if [[ -d $(build_dir) ]] && [[ $$(date --reference=$(build_dir) +%s) -le $$(date --date='-$(build_ttl)' +%s) ]]; then
		rm --recursive --force $(build_dir)
	fi
	if [[ ! -d $(build_dir) ]]; then
		git clone --depth=1 $(build_repo) $(build_dir)
	fi
	$(MAKE) $(MFLAGS) --no-print-directory --makefile=$(build_dir)/main.mk $(if $(MAKECMDGOALS),$@)
