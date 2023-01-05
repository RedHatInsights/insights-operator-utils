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

echo -e "${BLUE}Running formatting tool for Go source code${NC}"

if [ -n "$(gofmt -l .)" ]
then
    echo -e "${RED_BG}[FAIL]${NC} Go code is not formatted:"
    gofmt -d .
    exit 1
else
    echo -e "${GREEN_BG}[OK]${NC} Go code is formatted"
    exit 0
fi
