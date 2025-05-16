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
func UseBrokerConfig(brokerCfg *kafka.BrokerConfiguration, loadedConfig *api.AppConfig) {
	if loadedConfig.Kafka != nil && len(loadedConfig.Kafka.Brokers) > 0 {
		brokerCfg.Addresses = ""
		for _, broker := range loadedConfig.Kafka.Brokers {
			if broker.Port != nil {
				brokerCfg.Addresses += fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port) + ","
			} else {
				brokerCfg.Addresses += broker.Hostname + ","
			}
		}
		// remove the extra comma
		brokerCfg.Addresses = brokerCfg.Addresses[:len(brokerCfg.Addresses)-1]
		// SSL config
		clowderCfg := loadedConfig.Kafka.Brokers[0]
		if clowderCfg.Authtype != nil {
			fmt.Println("kafka is configured to use authentication")
			if clowderCfg.Sasl != nil {
				// we are trusting that these values are set and
				// dereferencing the pointers without any check...
				brokerCfg.SaslUsername = *clowderCfg.Sasl.Username
				brokerCfg.SaslPassword = *clowderCfg.Sasl.Password
				brokerCfg.SaslMechanism = *clowderCfg.Sasl.SaslMechanism
				brokerCfg.SecurityProtocol = *clowderCfg.SecurityProtocol

				if caPath, err := loadedConfig.KafkaCa(clowderCfg); err == nil {
					brokerCfg.CertPath = caPath
				}
			} else {
				fmt.Println(noSaslConfig)
			}
		}
	} else {
		fmt.Println(noBrokerConfig)
	}
}

// UseClowderTopics tries to replace the configured topic with the corresponding
// topic loaded by Clowder
func UseClowderTopics(brokerCfg *kafka.BrokerConfiguration, kafkaTopics map[string]api.TopicConfig) {
	if clowderTopic, ok := kafkaTopics[brokerCfg.Topic]; ok {
		brokerCfg.Topic = clowderTopic.Name
	} else {
		fmt.Printf(noTopicMapping, brokerCfg.Topic)
	}
}
