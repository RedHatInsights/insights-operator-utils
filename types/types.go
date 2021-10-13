// Copyright 2020, 2021  Red Hat, Inc
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

// Package types contains declaration of various data types (usually structures)
// used elsewhere in the aggregator code.
package types

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/types/types.html

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// OrgID represents organization ID
type OrgID uint32

// ClusterName represents name of cluster in format c8590f31-e97e-4b85-b506-c45ce1911a12
type ClusterName string

// UserID represents type for user id
type UserID string

// ClusterReport represents cluster report
type ClusterReport string

// Timestamp represents any timestamp in a form gathered from database
// TODO: need to be improved
type Timestamp string

// UserVote is a type for user's vote
type UserVote int

// RequestID is used to store the request ID supplied in input Kafka records as
// a unique identifier of payloads. Empty string represents a missing request ID.
type RequestID string

// RuleID represents type for rule id
type RuleID string

// RuleFQDN represents type for rule FQDN (rule module)
type RuleFQDN string

// RuleSelector represents component + error key
type RuleSelector string

// Component represent name of component (of rule)
type Component string

// RuleOnReport represents a single (hit) rule of the string encoded report
type RuleOnReport struct {
	Module          RuleID      `json:"component"`
	ErrorKey        ErrorKey    `json:"key"`
	UserVote        UserVote    `json:"user_vote"`
	Disabled        bool        `json:"disabled"`
	DisableFeedback string      `json:"disable_feedback"`
	DisabledAt      Timestamp   `json:"disabled_at"`
	TemplateData    interface{} `json:"details"`
}

// ReportRules is a helper struct for easy JSON unmarshalling of string encoded report
type ReportRules struct {
	HitRules     []RuleOnReport `json:"reports"`
	SkippedRules []RuleOnReport `json:"skips"`
	PassedRules  []RuleOnReport `json:"pass"`
	TotalCount   int
}

// ReportResponse represents the response of /report endpoint
type ReportResponse struct {
	Meta   ReportResponseMeta `json:"meta"`
	Report []RuleOnReport     `json:"reports"`
}

// ReportResponseMeta contains metadata about the report
type ReportResponseMeta struct {
	Count         int       `json:"count"`
	LastCheckedAt Timestamp `json:"last_checked_at"`
}

// RuleContentResponse represents a single rule in the response of /report endpoint
type RuleContentResponse struct {
	CreatedAt    string      `json:"created_at"`
	Description  string      `json:"description"`
	ErrorKey     string      `json:"-"`
	Generic      string      `json:"details"`
	Reason       string      `json:"reason"`
	Resolution   string      `json:"resolution"`
	TotalRisk    int         `json:"total_risk"`
	RiskOfChange int         `json:"risk_of_change"`
	RuleModule   RuleID      `json:"rule_id"`
	TemplateData interface{} `json:"extra_data"`
	Tags         []string    `json:"tags"`
	UserVote     UserVote    `json:"user_vote"`
	Disabled     bool        `json:"disabled"`
	Internal     bool        `json:"internal"`
}

// DisabledRuleResponse represents a single disabled rule displaying only identifying information
type DisabledRuleResponse struct {
	RuleModule  string `json:"rule_id"`
	Description string `json:"description"`
	Generic     string `json:"details"`
	DisabledAt  string `json:"disabled_at"`
}

// ErrorKey represents type for error key
type ErrorKey string

// Rule represents the content of rule table
type Rule struct {
	Module     RuleID `json:"module"`
	Name       string `json:"name"`
	Generic    string `json:"generic"`
	Summary    string `json:"summary"`
	Reason     string `json:"reason"`
	Resolution string `json:"resolution"`
	MoreInfo   string `json:"more_info"`
}

// RuleErrorKey represents the content of rule_error_key table
type RuleErrorKey struct {
	ErrorKey    ErrorKey  `json:"error_key"`
	RuleModule  RuleID    `json:"rule_module"`
	Condition   string    `json:"condition"`
	Description string    `json:"description"`
	Impact      int       `json:"impact"`
	Likelihood  int       `json:"likelihood"`
	PublishDate time.Time `json:"publish_date"`
	Active      bool      `json:"active"`
	Generic     string    `json:"generic"`
	Summary     string    `json:"summary"`
	Reason      string    `json:"reason"`
	Resolution  string    `json:"resolution"`
	MoreInfo    string    `json:"more_info"`
	Tags        []string  `json:"tags"`
}

// RuleWithContent represents a rule with content, basically the mix of rule and rule_error_key tables' content
type RuleWithContent struct {
	Module       RuleID    `json:"module"`
	Name         string    `json:"name"`
	Summary      string    `json:"summary"`
	Reason       string    `json:"reason"`
	Resolution   string    `json:"resolution"`
	MoreInfo     string    `json:"more_info"`
	ErrorKey     ErrorKey  `json:"error_key"`
	Condition    string    `json:"condition"`
	Description  string    `json:"description"`
	TotalRisk    int       `json:"total_risk"`
	RiskOfChange int       `json:"risk_of_change"`
	PublishDate  time.Time `json:"publish_date"`
	Active       bool      `json:"active"`
	Internal     bool      `json:"internal"`
	Generic      string    `json:"generic"`
	Tags         []string  `json:"tags"`
}

// ReportItem represents a single (hit) rule of the string encoded report
type ReportItem struct {
	Module       RuleID          `json:"component"`
	ErrorKey     ErrorKey        `json:"key"`
	TemplateData json.RawMessage `json:"details"`
}

// KafkaOffset type for kafka offset
type KafkaOffset int64

// DBDriver type for db driver enum
type DBDriver int

const (
	// DBDriverSQLite3 shows that db driver is sqlite
	DBDriverSQLite3 DBDriver = iota
	// DBDriverPostgres shows that db driver is postgres
	DBDriverPostgres
	// DBDriverGeneral general sql(used for mock now)
	DBDriverGeneral
)

const (
	// UserVoteDislike shows user's dislike
	UserVoteDislike UserVote = -1
	// UserVoteNone shows no vote from user
	UserVoteNone UserVote = 0
	// UserVoteLike shows user's like
	UserVoteLike UserVote = 1
)

// ValidationError validation error, for example when string is longer then expected
type ValidationError struct {
	ParamName  string
	ParamValue interface{}
	ErrString  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"Error during validating param '%v' with value '%v'. Error: '%v'",
		e.ParamName, e.ParamValue, e.ErrString,
	)
}

// ClusterListInRequest represents request body containing list of clusters
type ClusterListInRequest struct {
	Clusters []string `json:"clusters"`
}

// ClusterReports is a data structure containing list of clusters, list of
// errors and dictionary with results per cluster. This structure is used by
// aggregator to return more reports.
type ClusterReports struct {
	ClusterList []ClusterName                   `json:"clusters"`
	Errors      []ClusterName                   `json:"errors"`
	Reports     map[ClusterName]json.RawMessage `json:"reports"`
	GeneratedAt string                          `json:"generated_at"`
	Status      string                          `json:"status"`
}

// HittingClustersMetadata used to store metadata of clusters hit by a concrete rule
type HittingClustersMetadata struct {
	Count       int       `json:"count"`
	Component   Component `json:"component"`
	ErrorKey    ErrorKey  `json:"error_key"`
}

// HittingClustersData used to store data of clusters hit by a concrete rule
type HittingClustersData struct {
	Cluster  ClusterName `json:"cluster"`
	//Version  string      `json:"version"`
	GeneratedAt string   `json:"generated_at"`
}
// HittingClusters is a data structure containing list of clusters hit by a concrete rule
// hitting the given rule.
type HittingClusters struct {
	Metadata    HittingClustersMetadata `json:"meta"`
	ClusterList []HittingClustersData   `json:"data"`
}

//SchemaVersion is just a constant integer for now, max value 255. If we one day
//need more versions, better consider upgrading to semantic versioning.
type SchemaVersion uint8

// Acknowledgement represents user acknowledgement of given rule
type Acknowledgement struct {
	Acknowledged  bool   `json:"-"` // let's skip this one in responses
	Rule          string `json:"rule"`
	Justification string `json:"justification"`
	CreatedBy     string `json:"created_by"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// AcknowledgementsMetadata contains metadata about list of acknowledgements
type AcknowledgementsMetadata struct {
	Count int `json:"count"`
}

// AcknowledgementsResponse is structure returned to client in JSON
// serialization format
type AcknowledgementsResponse struct {
	Metadata AcknowledgementsMetadata `json:"meta"`
	Data     []Acknowledgement        `json:"data"`
}

// AcknowledgementJustification data structure represents body of request with
// specified justification of given acknowledgement
type AcknowledgementJustification struct {
	Value string `json:"justification"`
}

// AcknowledgementRuleSelectorJustification data structure represents body of
// request with specified rule selector and justification of given
// acknowledgement
type AcknowledgementRuleSelectorJustification struct {
	RuleSelector RuleSelector `json:"rule_id"`
	Value        string       `json:"justification"`
}

// RuleToggle is a type for user's vote
type RuleToggle int

// DisabledRule represents a record from rule_cluster_toggle
type DisabledRule struct {
	ClusterID  ClusterName
	RuleID     RuleID
	ErrorKey   ErrorKey
	Disabled   RuleToggle
	DisabledAt sql.NullTime
	EnabledAt  sql.NullTime
	UpdatedAt  sql.NullTime
}

// DisabledRuleReason represents a record from
// cluster_user_rule_disable_feedback table
type DisabledRuleReason struct {
	ClusterID ClusterName
	RuleID    RuleID
	ErrorKey  ErrorKey
	Message   string
	AddedAt   sql.NullTime
	UpdatedAt sql.NullTime
}

// SystemWideRuleDisable represents a record from rule_disable table
type SystemWideRuleDisable struct {
	OrgID         OrgID        `json:"org_id"`
	UserID        UserID       `json:"user_id"`
	RuleID        RuleID       `json:"rule_id"`
	ErrorKey      ErrorKey     `json:"error_key"`
	Justification string       `json:"justification"`
	CreatedAt     sql.NullTime `json:"created_at"`
	UpdatedAT     sql.NullTime `json:"updated_at"`
}

// ImpactedClustersCnt represents the number of clusters impacted by a rule
type ImpactedClustersCnt uint32

// RecommendationImpactedClusters is returned by aggregator for the purposes of /rule/ recommendation list endpoint
type RecommendationImpactedClusters map[RuleID]ImpactedClustersCnt

// RecommendationRow represents a single row in the recommendation table
type RecommendationRow struct {
	// RuleID is in "|" format
	RuleID    RuleID       `json:"rule_id"`
	RuleFQDN  RuleFQDN     `json:"rule_fqdn"`
	ErrorKey  ErrorKey     `json:"error_key"`
	OrgID     OrgID        `json:"org_id"`
	ClusterID ClusterName  `json:"cluster_id"`
	CreatedAt sql.NullTime `json:"created_at"`
}

// RecommendationListRow represents a single row retrieved from recommendation table
// for the purposes of Recommendations List (list of rules + number of impacted clusters)
type RecommendationListRow struct {
	RuleID              RuleID              `json:"rule_id"`
	ImpactedClustersCnt ImpactedClustersCnt `json:"impacted_clusters_cnt"`
}

// RuleRating represents the body request of request to the rating endpoint
type RuleRating struct {
	Rule   string   `json:"rule"`
	Rating UserVote `json:"rating"`
}

// RuleContentStatus type store information about rule content parsed and
// checked by Content Service
type RuleContentStatus struct {
	RuleType RuleType         `json:"type"`
	Loaded   bool             `json:"loaded"`
	Error    RuleParsingError `json:"error"`
}

// RuleType identifies whether the rule is external or internal one
// INFO: might be stored as a bool, but number of rule types might be enhanced
// later
type RuleType string

// RuleParsingError represents textual and human-readable form of (any) error
// occured during reading, parsing, and checking rule content in Content
// Service
type RuleParsingError string
