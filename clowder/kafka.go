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
	"fmt"

	"github.com/RedHatInsights/insights-operator-utils/kafka"
	api "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

// Common constants used for logging and error reporting
const (
	noOriginalBroker = "warning: no original broker configuration found; aborting"
	noBrokerConfig   = "warning: no broker configurations found in clowder config"
	noSaslConfig     = "warning: SASL configuration is missing"
	noTopicMapping   = "warning: no kafka mapping found for topic %s"
)

// UseBrokerConfig tries to replace parts of the BrokerConfiguration with the values
// loaded by Clowder. It expects brokerCfg to already be initialized
func UseBrokerConfig(brokerCfg kafka.BrokersConfig, loadedConfig *api.AppConfig) kafka.BrokersConfig {
	numBrokerConfigs := len(brokerCfg)
	if numBrokerConfigs == 0 {
		// if original brokers config is totally empty, do nothing.
		// this shouldn't happen, but we need to control this scenario
		// to avoid panics.
		fmt.Println(noOriginalBroker)
		return brokerCfg
	}
	if loadedConfig.Kafka != nil && len(loadedConfig.Kafka.Brokers) > 0 {
		numClowderBrokers := len(loadedConfig.Kafka.Brokers)
		// if original config has fewer brokers than clowder's, we append additional
		// brokerConfiguration items with topic, clientId, and group from existing
		// items, and the rest will be filled with data from clowder's brokers.
		// When appending, it's most probable that a new slice is returned due to
		// original capacity not being enough, which is why the brokerCfg slice is
		// returned
		for len(brokerCfg) < numClowderBrokers {
			brokerCfg = append(brokerCfg, &kafka.BrokerConfiguration{
				Topic:    (brokerCfg)[numBrokerConfigs-1].Topic,
				ClientID: (brokerCfg)[numBrokerConfigs-1].ClientID,
				Group:    (brokerCfg)[numBrokerConfigs-1].Group,
				// Since this will have data from Clowder's loadedConfig, always enable
				Enabled: true,
			})
		}
		for i, broker := range loadedConfig.Kafka.Brokers {
			if broker.Port != nil {
				(brokerCfg)[i].Address = fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port)
			} else {
				(brokerCfg)[i].Address = broker.Hostname
			}
			// SSL config
			if broker.Authtype != nil {
				fmt.Println("kafka is configured to use authentication")
				if broker.Sasl != nil {
					// we are trusting that these values are set and
					// dereferencing the pointers without any check...
					brokerCfg[i].SaslUsername = *broker.Sasl.Username
					brokerCfg[i].SaslPassword = *broker.Sasl.Password
					brokerCfg[i].SaslMechanism = *broker.Sasl.SaslMechanism
					brokerCfg[i].SecurityProtocol = *broker.SecurityProtocol

					if caPath, err := loadedConfig.KafkaCa(broker); err == nil {
						brokerCfg[i].CertPath = caPath
					}
				} else {
					fmt.Println(noSaslConfig)
				}
			}
		}
	} else {
		fmt.Println(noBrokerConfig)
	}
	return brokerCfg
}

// UseClowderTopics tries to replace the configured topic's name with the
// corresponding topic name loaded by Clowder, if any
func UseClowderTopics(brokersCfg kafka.BrokersConfig, kafkaTopics map[string]api.TopicConfig) {
	for _, cfg := range brokersCfg {
		if clowderTopic, ok := kafkaTopics[cfg.Topic]; ok {
			cfg.Topic = clowderTopic.Name
		} else {
			fmt.Printf(noTopicMapping, cfg.Topic)
		}
	}
}
