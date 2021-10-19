// Copyright 2020 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/types/content.html

// RuleContent wraps all the content available for a rule into a single structure.
type RuleContent struct {
	Plugin     RulePluginInfo                 `json:"plugin"`
	ErrorKeys  map[string]RuleErrorKeyContent `json:"error_keys"`
	Generic    string                         `json:"generic"`
	Summary    string                         `json:"summary"`
	Resolution string                         `json:"resolution"`
	MoreInfo   string                         `json:"more_info"`
	Reason     string                         `json:"reason"`
	HasReason  bool
}

// RulePluginInfo is a Go representation of the `plugin.yaml`
// file inside of the rule content directory.
type RulePluginInfo struct {
	Name         string `yaml:"name" json:"name"`
	NodeID       string `yaml:"node_id" json:"node_id"`
	ProductCode  string `yaml:"product_code" json:"product_code"`
	PythonModule string `yaml:"python_module" json:"python_module"`
}

// RuleErrorKeyContent wraps content of a single error key.
type RuleErrorKeyContent struct {
	Metadata   ErrorKeyMetadata `json:"metadata"`
	TotalRisk  int              `json:"total_risk"`
	Generic    string           `json:"generic"`
	Summary    string           `json:"summary"`
	Resolution string           `json:"resolution"`
	MoreInfo   string           `json:"more_info"`
	Reason     string           `json:"reason"`
	// DONTFIX has_reason until CCXDEV-5021
	HasReason bool
}

// ErrorKeyMetadata is a Go representation of the `metadata.yaml`
// file inside of an error key content directory.
type ErrorKeyMetadata struct {
	Description string   `yaml:"description" json:"description"`
	Impact      Impact   `yaml:"impact" json:"impact"`
	Likelihood  int      `yaml:"likelihood" json:"likelihood"`
	PublishDate string   `yaml:"publish_date" json:"publish_date"`
	Status      string   `yaml:"status" json:"status"`
	Tags        []string `yaml:"tags" json:"tags"`
}

// Impact is contained in ErrorKeyMetadata
type Impact struct {
	Name   string `yaml:"name" json:"name"`
	Impact int    `yaml:"impact" json:"impact"`
}

// RuleContentDirectory contains content for all available rules in a directory.
type RuleContentDirectory struct {
	Config GlobalRuleConfig
	Rules  map[string]RuleContent
}

// GlobalRuleConfig represents the file that contains
// metadata globally applicable to any/all rule content.
type GlobalRuleConfig struct {
	Impact map[string]int `yaml:"impact" json:"impact"`
}
