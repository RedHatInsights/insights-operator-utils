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

func UseBrokerConfig(brokerCfg *kafka.BrokerConfiguration, loadedConfig *api.AppConfig) {
	// make sure broker(s) are configured in Clowder
	if loadedConfig.Kafka != nil && len(loadedConfig.Kafka.Brokers) > 0 {
		broker := loadedConfig.Kafka.Brokers[0]
		// port can be empty in api, so taking it into account
		if broker.Port != nil {
			brokerCfg.Address = fmt.Sprintf("%s:%d", broker.Hostname, *broker.Port)
		} else {
			brokerCfg.Address = broker.Hostname
		}

		// SSL config
		if broker.Authtype != nil {
			fmt.Println("kafka is configured to use authentication")
			if broker.Sasl != nil {
				// we are trusting that these values are set and
				// dereferencing the pointers without any check...
				brokerCfg.SaslUsername = *broker.Sasl.Username
				brokerCfg.SaslPassword = *broker.Sasl.Password
				brokerCfg.SaslMechanism = *broker.Sasl.SaslMechanism
				brokerCfg.SecurityProtocol = *broker.SecurityProtocol

				if caPath, err := loadedConfig.KafkaCa(broker); err == nil {
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
func UseClowderTopics(configuration *kafka.BrokerConfiguration, kafkaTopics map[string]api.TopicConfig) {
	// Get the correct topic name from clowder mapping if available
	if clowderTopic, ok := kafkaTopics[configuration.Topic]; ok {
		configuration.Topic = clowderTopic.Name
	} else {
		fmt.Printf(noTopicMapping, configuration.Topic)
	}
}
