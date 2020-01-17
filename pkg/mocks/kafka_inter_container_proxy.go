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

package mocks

import (
	"context"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/opencord/voltha-lib-go/v2/pkg/kafka"
)

type InvokeRpcArgs struct {
	Rpc             string
	ToTopic         *kafka.Topic
	ReplyToTopic    *kafka.Topic
	WaitForResponse bool
	Key             string
	ParentDeviceId  string
	KvArgs          map[int]interface{}
}

type InvokeRpcSpy struct {
	CallCount int
	Calls     map[int]InvokeRpcArgs
}

type MockKafkaIcProxy struct {
	InvokeRpcSpy InvokeRpcSpy
}

func (s *MockKafkaIcProxy) Start() error                        { return nil }
func (s *MockKafkaIcProxy) DeleteTopic(topic kafka.Topic) error { return nil }
func (s *MockKafkaIcProxy) DeviceDiscovered(deviceId string, deviceType string, parentId string, publisher string) error {
	return nil
}
func (s *MockKafkaIcProxy) Stop() {}
func (s *MockKafkaIcProxy) InvokeRPC(ctx context.Context, rpc string, toTopic *kafka.Topic, replyToTopic *kafka.Topic, waitForResponse bool, key string, kvArgs ...*kafka.KVArg) (bool, *any.Any) {
	s.InvokeRpcSpy.CallCount++

	args := make(map[int]interface{}, 4)
	for k, v := range kvArgs {
		args[k] = v
	}

	s.InvokeRpcSpy.Calls[s.InvokeRpcSpy.CallCount] = InvokeRpcArgs{
		Rpc:             rpc,
		ToTopic:         toTopic,
		ReplyToTopic:    replyToTopic,
		WaitForResponse: waitForResponse,
		Key:             key,
		KvArgs:          args,
	}
	return false, &any.Any{}
}
func (s *MockKafkaIcProxy) SubscribeWithRequestHandlerInterface(topic kafka.Topic, handler interface{}) error {
	return nil
}
func (s *MockKafkaIcProxy) SubscribeWithDefaultRequestHandler(topic kafka.Topic, initialOffset int64) error {
	return nil
}
func (s *MockKafkaIcProxy) UnSubscribeFromRequestHandler(topic kafka.Topic) error { return nil }
