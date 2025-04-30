// Copyright 2023 Red Hat, Inc
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

// Package services contains interface implementations to other
// services that are called from Smart Proxy.
package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/redis"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

var (
	ctx                        context.Context
	defaultRedisAddress        = "loclahost:6379"
	defaultRedisDatabase       = 0
	defaultRedisUsername       = "default"
	defaultRedisPassword       = "psw"
	defaultRedisTimeoutSeconds = 30
)

// set default conf
func init() {
	ctx = context.Background()
}

func getMockRedis() (
	mockClient redis.Client, mockServer redismock.ClientMock,
) {
	client, mockServer := redismock.NewClientMock()
	mockClient = redis.Client{
		Connection: client,
	}
	return
}

func redisExpectationsMet(t *testing.T, mock redismock.ClientMock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestCreateRedisClientOK(t *testing.T) {
	client, err := redis.CreateRedisClient(
		defaultRedisAddress, defaultRedisDatabase, defaultRedisUsername, defaultRedisPassword, defaultRedisTimeoutSeconds,
	)
	assert.NoError(t, err)

	options := client.Options()
	assert.NoError(t, err)
	assert.Equal(t, options.Addr, defaultRedisAddress)
	assert.Equal(t, options.DB, defaultRedisDatabase)
	assert.Equal(t, options.Password, defaultRedisPassword)
	assert.Equal(t, options.ReadTimeout, time.Duration(defaultRedisTimeoutSeconds)*time.Second)
}

func TestCreateRedisClientBadAddress(t *testing.T) {
	// empty address
	client, err := redis.CreateRedisClient(
		"", defaultRedisDatabase, defaultRedisUsername, defaultRedisPassword, defaultRedisTimeoutSeconds,
	)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestCreateRedisClientDBIndexOutOfRange(t *testing.T) {
	// Redis supports "only" 16 different databases with indices 0-15
	client, err := redis.CreateRedisClient(
		defaultRedisAddress, 16, defaultRedisUsername, defaultRedisPassword, defaultRedisTimeoutSeconds,
	)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestRedisHealthCheckOK(t *testing.T) {
	client, server := getMockRedis()

	server.ExpectPing().SetVal("PONG")

	err := client.HealthCheck()
	assert.NoError(t, err)

	redisExpectationsMet(t, server)
}

func TestRedisHealthCheckError(t *testing.T) {
	client, server := getMockRedis()

	server.ExpectPing().SetErr(errors.New("mock error"))

	err := client.HealthCheck()
	assert.Error(t, err)

	redisExpectationsMet(t, server)
}

func TestRedisHealthCheckBadResponse(t *testing.T) {
	client, server := getMockRedis()

	// cover 2nd condition
	server.ExpectPing().SetVal("ka-boom")

	err := client.HealthCheck()
	assert.Error(t, err)

	redisExpectationsMet(t, server)
}
