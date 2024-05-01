# -*- makefile -*-
# -----------------------------------------------------------------------
# Copyright 2022-2024 Open Networking Foundation (ONF) and the ONF Contributors
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
#
# SPDX-FileCopyrightText: 2022-2023 Open Networking Foundation (ONF) and the ONF Contributors
# SPDX-License-Identifier: Apache-2.0
# -----------------------------------------------------------------------
# https://gerrit.opencord.org/plugins/gitiles/onf-make
# ONF.makefile.version = 1.1
# -----------------------------------------------------------------------

$(if $(DEBUG),$(warning ENTER))

## -----------------------------------------------------------------------
## -----------------------------------------------------------------------
help ::
	@echo
	@echo "[LINT]"

## moved to config.mk
## Disable python linting in bulk ?
# ifdef NO-LINT-PYTHON
#   NO-LINT-FLAKE8 := true
#   NO-LINT-PYLINT := true
# endif

include $(legacy-mk)/lint/groovy.mk
include $(legacy-mk)/lint/jjb.mk
include $(legacy-mk)/lint/json.mk
include $(legacy-mk)/lint/license/include.mk
include $(legacy-mk)/lint/makefile.mk
include $(legacy-mk)/lint/python/include.mk
include $(legacy-mk)/lint/shellcheck/include.mk
include $(legacy-mk)/lint/tox/include.mk
include $(legacy-mk)/lint/yaml/include.mk

include $(legacy-mk)/lint/help.mk

$(if $(DEBUG),$(warning LEAVE))

# [EOF]
