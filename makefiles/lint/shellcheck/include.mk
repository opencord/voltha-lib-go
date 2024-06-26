# -*- makefile -*-
# -----------------------------------------------------------------------
# Copyright 2017-2024 Open Networking Foundation (ONF) and the ONF Contributors
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

ifndef mk-include--onf-lint-shellcheck # single-include guard macro

$(if $(DEBUG),$(warning ENTER))

##--------------------##
##---]  INCLUDES  [---##
##--------------------##
# include $(legacy-mk)/lint/shellcheck/help.mk
include $(legacy-mk)/lint/shellcheck/find_utils.mk

# Standard lint-yaml targets
include $(legacy-mk)/lint/shellcheck/shellcheck.mk

mk-include--onf-lint-shellcheck := true#        # Flag to inhibit re-including

$(if $(DEBUG),$(warning LEAVE))

endif # mk-include--onf-lint-shellcheck

# [EOF]
