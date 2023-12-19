#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

# Copyright (c) 2023  The Go-Curses Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#: uncomment to echo instead of execute
#CMD=echo

-include .env
#export

BIN_NAME := go-mod-update
UNTAGGED_VERSION := v0.1.1
UNTAGGED_COMMIT := trunk

SHELL := /bin/bash
RUN_ARGS := --help
LOG_LEVEL := debug

GO_ENJIN_PKG := nil
BE_LOCAL_PATH := nil

GOPKG_KEYS ?= CDK CTK CLR CLC CLS CLM

CDK_GO_PACKAGE ?= github.com/go-curses/cdk
CDK_LOCAL_PATH ?= ../cdk

CTK_GO_PACKAGE ?= github.com/go-curses/ctk
CTK_LOCAL_PATH ?= ../ctk

CLR_GO_PACKAGE ?= github.com/go-curses/corelibs/run
CLR_LOCAL_PATH ?= ../corelibs/run

CLC_GO_PACKAGE ?= github.com/go-curses/corelibs/chdirs
CLC_LOCAL_PATH ?= ../corelibs/chdirs

CLS_GO_PACKAGE ?= github.com/go-curses/corelibs/spinner
CLS_LOCAL_PATH ?= ../corelibs/spinner

CLM_GO_PACKAGE ?= github.com/go-curses/corelibs/maps
CLM_LOCAL_PATH ?= ../corelibs/maps

CLEAN_FILES     ?= ${BIN_NAME} ${BIN_NAME}.*.* coverage.out pprof.*
DISTCLEAN_FILES ?=
REALCLEAN_FILES ?=

BUILD_VERSION_VAR := main.APP_VERSION
BUILD_RELEASE_VAR := main.APP_RELEASE

SRC_CMD_PATH := ./cmd/go-mod-update

include Golang.cmd.mk
include Golang.def.mk
include Golang.cdk.mk
