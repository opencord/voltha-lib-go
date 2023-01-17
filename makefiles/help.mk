# -*- makefile -*-
# -----------------------------------------------------------------------
# Copyright 2022-2023 Open Networking Foundation (ONF) and the ONF Contributors (ONF) and the ONF Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -----------------------------------------------------------------------

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
.PHONY: help
help ::
	@echo "Usage: $(MAKE) [options] [target] ..."
	@echo
	@echo "  build                : Build the library"
	@echo "  mod-update           : Update go.mod and the vendor directory"
	@echo "  test                 : Generate reports for all go tests"
	@echo "  local-protos         : Local development target"
ifdef VERBOSE
	@echo "   LOCAL_PROTOS="
endif
	@echo
	@echo "[CLEAN]"
	@echo "  clean                : Remove files created by the build"
	@echo "  distclean            : Remove build and testing artifacts and reports"
	@echo
	@echo "[LINT]"
	@echo "  lint                 : Shorthand for lint-style & lint-sanity"
	@echo "  lint-mod             : Verify the integrity of the 'mod' files"
	@echo "  sca                  : "

# [EOF]
