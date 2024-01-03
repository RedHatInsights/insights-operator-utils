package clowder

import (
	"github.com/RedHatInsights/insights-operator-utils/postgres"
	api "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func UseDBConfig(storageCfg *postgres.StorageConfiguration, loadedConfig *api.AppConfig) {
	storageCfg.PGDBName = loadedConfig.Database.Name
	storageCfg.PGHost = loadedConfig.Database.Hostname
	storageCfg.PGPort = loadedConfig.Database.Port
	storageCfg.PGUsername = loadedConfig.Database.Username
	storageCfg.PGPassword = loadedConfig.Database.Password
}
