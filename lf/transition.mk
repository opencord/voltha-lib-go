# -*- makefile -*-
# -----------------------------------------------------------------------
# Copyright 2024 Open Networking Foundation Contributors
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
# SPDX-FileCopyrightText: 2024 Open Networking Foundation Contributors
# SPDX-License-Identifier: Apache-2.0
# -----------------------------------------------------------------------

$(if $(DEBUG),$(warning ENTER))

##-------------------##
##---]  GLOBALS  [---##
##-------------------##
lf-sbx-root   := $(abspath $(lastword $(MAKEFILE_LIST)))
lf-sbx-root   := $(subst /lf/transition.mk,$(null),$(lf-sbx-root))

legacy-mk   ?= $(lf-sbx-root)/makefiles
onf-mk-dir  ?= $(lf-sbx-root)/lf/onf-make/makefiles

sandbox-root := $(lf-sbx-root)

# set default shell options
SHELL = bash -e -o pipefail

# dependency of virtualenv::sterile
clean ::

##--------------------##
##---]  INCLUDES  [---##
##--------------------##
include $(lf-sbx-root)/lf/config.mk
include $(legacy-mk)/include.mk

# [EOF]

ifdef paste-into-repository-makefile

$(if $(findstring disabled-joey,$(USER)),\
   $(eval USE_LF_MK := 1)) # special snowflake

##--------------------##
##---]  INCLUDES  [---##
##--------------------##
ifdef USE_LF_MK
  include lf/include.mk
else
  include lf/transition.mk
endif # ifdef USE_LF_MK

endif # paste-into-repository-makefile
