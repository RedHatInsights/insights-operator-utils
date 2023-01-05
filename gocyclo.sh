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

echo -e "${BLUE}Finding functions and methods with high cyclomatic complexity${NC}"


if ! [ -x "$(command -v gocyclo)" ]
then
    echo -e "${BLUE}Installing gocyclo${NC}"
    GO111MODULE=off go get github.com/fzipp/gocyclo/cmd/gocyclo
fi

if ! gocyclo -over 13 -avg .
then
    echo -e "${RED_BG}[FAIL]${NC} Functions/methods with high cyclomatic complexity detected"
    exit 1
else
    echo -e "${GREEN_BG}[OK]${NC} All functions and methods have reasonable cyclomatic complexity"
    exit 0
fi
