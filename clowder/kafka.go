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
	noBrokerConfig = "warning: no broker configurations found in clowder config"
	noSaslConfig   = "warning: SASL configuration is missing"
	noTopicMapping = "warning: no kafka mapping found for topic %s"
)

// UseBrokerConfig tries to replace parts of the BrokerConfiguration with the values
// loaded by Clowder
func UseBrokerConfig(brokerCfg *kafka.MultiBrokerConfiguration, loadedConfig *api.AppConfig) {
	if loadedConfig.Kafka != nil && len(loadedConfig.Kafka.Brokers) > 0 {
		brokerCfg.Addresses = make([]string, len(loadedConfig.Kafka.Brokers))
		brokerCfg.SASLConfigs = make([]kafka.SASLConfiguration, len(loadedConfig.Kafka.Brokers))
		for i, broker := range loadedConfig.Kafka.Brokers {
			if broker.Port != nil {
				brokerCfg.Addresses[i] = fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port)
			} else {
				brokerCfg.Addresses[i] = broker.Hostname
			}
			// SSL config
			if broker.Authtype != nil {
				fmt.Println("kafka is configured to use authentication")
				if broker.Sasl != nil {
					// we are trusting that these values are set and
					// dereferencing the pointers without any check...
					brokerCfg.SASLConfigs[i].SaslUsername = *broker.Sasl.Username
					brokerCfg.SASLConfigs[i].SaslPassword = *broker.Sasl.Password
					brokerCfg.SASLConfigs[i].SaslMechanism = *broker.Sasl.SaslMechanism
					brokerCfg.SASLConfigs[i].SecurityProtocol = *broker.SecurityProtocol

					if caPath, err := loadedConfig.KafkaCa(broker); err == nil {
						brokerCfg.SASLConfigs[i].CertPath = caPath
					}
				} else {
					fmt.Println(noSaslConfig)
				}
			}
		}
	} else {
		fmt.Println(noBrokerConfig)
	}
}

// UseClowderTopics tries to replace the configured topic with the corresponding
// topic loaded by Clowder
func UseClowderTopics(brokerCfg interface{}, kafkaTopics map[string]api.TopicConfig) {
	switch cfg := brokerCfg.(type) {
	case *kafka.SingleBrokerConfiguration:
		if clowderTopic, ok := kafkaTopics[cfg.Topic]; ok {
			cfg.Topic = clowderTopic.Name
		} else {
			fmt.Printf(noTopicMapping, cfg.Topic)
		}
	case *kafka.MultiBrokerConfiguration:
		if clowderTopic, ok := kafkaTopics[cfg.Topic]; ok {
			cfg.Topic = clowderTopic.Name
		} else {
			fmt.Printf(noTopicMapping, cfg.Topic)
		}
	default:
		fmt.Printf("Unknown Broker configuration type")
	}
}
