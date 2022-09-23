#!/bin/bash

# Copyright 2020, 2021 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BLUE=$(tput setaf 4)
RED_BG=$(tput setab 1)
GREEN_BG=$(tput setab 2)
NC=$(tput sgr0) # No Color

echo -e "${BLUE}Running linter for Go source code${NC}"

cd "$(dirname "$0")" || exit

if ! [ -x "$(command -v golint)" ]
then
    echo -e "${BLUE}Installing golint${NC}"
    GO111MODULE=off go get golang.org/x/lint/golint 2> /dev/null
fi

# shellcheck disable=SC2046
if golint $(go list ./...) |
    grep -v ALL_CAPS |
    grep .
then
    echo -e "${RED_BG}[FAIL]${NC} Linter have found some issues in Go source code"
    exit 1
else
    echo -e "${GREEN_BG}[OK]${NC} No issues has been found"
    exit 0
fi
