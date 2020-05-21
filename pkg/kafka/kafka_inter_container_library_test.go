/*
 * Copyright 2018-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package kafka

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultKafkaProxy(t *testing.T) {
	actualResult := newInterContainerProxy(context.Background())
	assert.Equal(t, actualResult.kafkaHost, DefaultKafkaHost)
	assert.Equal(t, actualResult.kafkaPort, DefaultKafkaPort)
	assert.Equal(t, actualResult.defaultRequestHandlerInterface, interface{}(nil))
}

func TestKafkaProxyOptionHost(t *testing.T) {
	actualResult := newInterContainerProxy(context.Background(), InterContainerHost(context.Background(), "10.20.30.40"))
	assert.Equal(t, actualResult.kafkaHost, "10.20.30.40")
	assert.Equal(t, actualResult.kafkaPort, DefaultKafkaPort)
	assert.Equal(t, actualResult.defaultRequestHandlerInterface, interface{}(nil))
}

func TestKafkaProxyOptionPort(t *testing.T) {
	actualResult := newInterContainerProxy(context.Background(), InterContainerPort(context.Background(), 1020))
	assert.Equal(t, actualResult.kafkaHost, DefaultKafkaHost)
	assert.Equal(t, actualResult.kafkaPort, 1020)
	assert.Equal(t, actualResult.defaultRequestHandlerInterface, interface{}(nil))
}

func TestKafkaProxyOptionTopic(t *testing.T) {
	actualResult := newInterContainerProxy(context.Background(), DefaultTopic(context.Background(), &Topic{Name: "Adapter"}))
	assert.Equal(t, actualResult.kafkaHost, DefaultKafkaHost)
	assert.Equal(t, actualResult.kafkaPort, DefaultKafkaPort)
	assert.Equal(t, actualResult.defaultRequestHandlerInterface, interface{}(nil))
	assert.Equal(t, actualResult.defaultTopic.Name, "Adapter")
}

type myInterface struct {
}

func TestKafkaProxyOptionTargetInterface(t *testing.T) {
	var m *myInterface
	actualResult := newInterContainerProxy(context.Background(), RequestHandlerInterface(context.Background(), m))
	assert.Equal(t, actualResult.kafkaHost, DefaultKafkaHost)
	assert.Equal(t, actualResult.kafkaPort, DefaultKafkaPort)
	assert.Equal(t, actualResult.defaultRequestHandlerInterface, m)
}

func TestKafkaProxyChangeAllOptions(t *testing.T) {
	var m *myInterface
	actualResult := newInterContainerProxy(
		context.Background(),
		InterContainerHost(context.Background(), "10.20.30.40"),
		InterContainerPort(context.Background(), 1020),
		DefaultTopic(context.Background(), &Topic{Name: "Adapter"}),
		RequestHandlerInterface(context.Background(), m))
	assert.Equal(t, actualResult.kafkaHost, "10.20.30.40")
	assert.Equal(t, actualResult.kafkaPort, 1020)
	assert.Equal(t, actualResult.defaultRequestHandlerInterface, m)
	assert.Equal(t, actualResult.defaultTopic.Name, "Adapter")
}

func TestKafkaProxyEnableLivenessChannel(t *testing.T) {
	var m *myInterface

	// Note: This doesn't actually start the client
	client := NewSaramaClient(context.Background())

	probe := newInterContainerProxy(
		context.Background(),
		InterContainerHost(context.Background(), "10.20.30.40"),
		InterContainerPort(context.Background(), 1020),
		DefaultTopic(context.Background(), &Topic{Name: "Adapter"}),
		RequestHandlerInterface(context.Background(), m),
		MsgClient(context.Background(), client),
	)

	ch := probe.EnableLivenessChannel(context.Background(), true)

	// The channel should have one "true" message on it
	assert.NotEmpty(t, ch)

	select {
	case stuff := <-ch:
		assert.True(t, stuff)
	default:
		t.Error("Failed to read from the channel")
	}
}
