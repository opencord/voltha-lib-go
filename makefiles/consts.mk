# -*- makefile -*-
# -----------------------------------------------------------------------
# Copyright 2022 Open Networking Foundation (ONF) and the ONF Contributors
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
# SPDX-FileCopyrightText: 2022 Open Networking Foundation (ONF) and the ONF Contributors
# SPDX-License-Identifier: Apache-2.0
# -----------------------------------------------------------------------

$(if $(DEBUG),$(warning ENTER))

null        :=#
space       := $(null) $(null)
dot         ?= .

HIDE        ?= @

env-clean   = /usr/bin/env --ignore-environment
xargs-n1    := xargs -0 -t -n1 --no-run-if-empty

## -----------------------------------------------------------------------
## Not recommended but support (-u)ndef-less shell for pyenv activate
## TODO: declare a pyenv shell
## -----------------------------------------------------------------------
have-shell-bash := $(filter bash,$(subst /,$(space),$(SHELL)))
$(if $(have-shell-bash),$(null),\
  $(eval export SHELL := /bin/bash -euo pipefail))

shell-pyenv := bash -eo pipefail

$(if $(DEBUG),$(warning LEAVE))

# [EOF]
