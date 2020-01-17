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
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/opencord/voltha-lib-go/v3/pkg/kafka"
	ic "github.com/opencord/voltha-protos/v3/go/inter_container"
	"github.com/opencord/voltha-protos/v3/go/voltha"
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

type FailReason int

const (
	Timeout FailReason = iota + 1
	UnmarshalError
)

func (r FailReason) String() string {
	return [...]string{"Timeout", "UnmarshalError"}[r]
}

type InvokeRpcSpy struct {
	CallCount int
	Calls     map[int]InvokeRpcArgs
	Fail      FailReason // timeout, error
}

type MockKafkaICProxy struct {
	InvokeRpcSpy InvokeRpcSpy
}

func (s *MockKafkaICProxy) Start() error                        { return nil }
func (s *MockKafkaICProxy) DeleteTopic(topic kafka.Topic) error { return nil }
func (s *MockKafkaICProxy) DeviceDiscovered(deviceId string, deviceType string, parentId string, publisher string) error {
	return nil
}
func (s *MockKafkaICProxy) Stop() {}
func (s *MockKafkaICProxy) InvokeRPC(ctx context.Context, rpc string, toTopic *kafka.Topic, replyToTopic *kafka.Topic, waitForResponse bool, key string, kvArgs ...*kafka.KVArg) (bool, *any.Any) {
	s.InvokeRpcSpy.CallCount++

	success := true

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

	device := &voltha.Device{
		Id: "testDevice",
	}
	response, _ := ptypes.MarshalAny(device)

	if s.InvokeRpcSpy.Fail == Timeout {

		success = false

		// TODO once InvokeRPC is fixed to return an error code, add it here
		err := &ic.Error{Reason: "context deadline exceeded"}
		response, _ = ptypes.MarshalAny(err)
	} else if s.InvokeRpcSpy.Fail == UnmarshalError {
		res := &voltha.LogicalDevice{
			Id: "testLogicalDevice",
		}
		response, _ = ptypes.MarshalAny(res)
	}

	return success, response
}
func (s *MockKafkaICProxy) SubscribeWithRequestHandlerInterface(topic kafka.Topic, handler interface{}) error {
	return nil
}
func (s *MockKafkaICProxy) SubscribeWithDefaultRequestHandler(topic kafka.Topic, initialOffset int64) error {
	return nil
}
func (s *MockKafkaICProxy) UnSubscribeFromRequestHandler(topic kafka.Topic) error { return nil }
