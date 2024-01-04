// Copyright 2024 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clowder

import (
	"github.com/RedHatInsights/insights-operator-utils/postgres"
	api "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

// UseDBConfig tries to replace the StorageConfiguration parameters with the
// values loaded by Clowder
func UseDBConfig(storageCfg *postgres.StorageConfiguration, loadedConfig *api.AppConfig) {
	storageCfg.PGDBName = loadedConfig.Database.Name
	storageCfg.PGHost = loadedConfig.Database.Hostname
	storageCfg.PGPort = loadedConfig.Database.Port
	storageCfg.PGUsername = loadedConfig.Database.Username
	storageCfg.PGPassword = loadedConfig.Database.Password
}
