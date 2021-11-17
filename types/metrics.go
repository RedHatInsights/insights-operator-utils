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

// MetricsConfiguration holds metrics related configuration
type MetricsConfiguration struct {
	Job              string `mapstructure:"job_name" toml:"job_name"`
	GatewayURL       string `mapstructure:"gateway_url" toml:"gateway_url"`
	GatewayAuthToken string `mapstructure:"gateway_auth_token" toml:"gateway_auth_token"`
	TimeBetweenPush  int    `mapstructure:"time_between_push" toml:"time_between_push"`
}
