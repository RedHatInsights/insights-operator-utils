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

// Package redis contains shared functionality related to Redis
package redis

import (
	"context"
	"errors"
	"time"

	redisV9 "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	redisExecutionFailedMsg = "unexpected response from Redis server"
)

// Client is an implementation of the Redis client for Redis server
type Client struct {
	Connection *redisV9.Client
}

// CreateRedisClient creates a Redis V9 client, it has explicit checks for config params
// because the go-redis package lets us create a client with incorrect params
// so errors are only found during subsequent command executions
func CreateRedisClient(
	address string,
	databaseIndex int,
	username string,
	password string,
	timeoutSeconds int,
) (*redisV9.Client, error) {
	if address == "" {
		err := errors.New("Redis server address must not be empty")
		log.Error().Err(err)
		return nil, err
	}

	if databaseIndex < 0 || databaseIndex > 15 {
		err := errors.New("Redis selected database must be a value in the range 0-15")
		log.Error().Err(err)
		return nil, err
	}

	log.Info().Msgf("creating redis client. endpoint %v, selected DB %d, timeout seconds %d",
		address, databaseIndex, timeoutSeconds,
	)

	// DB is configurable in case we want to change data structure
	c := redisV9.NewClient(&redisV9.Options{
		Addr:        address,
		DB:          databaseIndex,
		Username:    username,
		Password:    password,
		ReadTimeout: time.Duration(timeoutSeconds) * time.Second,
	})

	return c, nil
}

// HealthCheck executes PING command to check for liveness status of Redis server
func (redis *Client) HealthCheck() (err error) {
	ctx := context.Background()

	// .Result() gets value and error of executed command at once
	res, err := redis.Connection.Ping(ctx).Result()
	if err != nil || res != "PONG" {
		log.Error().Err(err).Msg("Redis PING command failed")
		return errors.New(redisExecutionFailedMsg)
	}

	return
}
