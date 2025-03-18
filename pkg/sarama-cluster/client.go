/*
 * Copyright 2018-2024 Open Networking Foundation (ONF) and the ONF Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cluster

import (
	"errors"
	"sync/atomic"

	"github.com/Shopify/sarama"
)

var errClientInUse = errors.New("cluster: client is already used by another consumer")

// Client is a group client
type Client struct {
	sarama.Client
	config Config

	inUse uint32
}

// NewClient creates a new client instance
func NewClient(addrs []string, config *Config) (*Client, error) {
	if config == nil {
		config = NewConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	client, err := sarama.NewClient(addrs, &config.Config)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client, config: *config}, nil
}

// ClusterConfig returns the cluster configuration.
func (c *Client) ClusterConfig() *Config {
	cfg := c.config
	return &cfg
}

func (c *Client) claim() bool {
	return atomic.CompareAndSwapUint32(&c.inUse, 0, 1)
}

func (c *Client) release() {
	atomic.CompareAndSwapUint32(&c.inUse, 1, 0)
}
