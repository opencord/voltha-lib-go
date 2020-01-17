/*
 * Copyright 2019-present Open Networking Foundation
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
package common

import (
	"context"
	"github.com/golang/protobuf/ptypes/any"
	adapterIf "github.com/opencord/voltha-lib-go/v2/pkg/adapters/adapterif"
	ic "github.com/opencord/voltha-protos/v2/go/inter_container"
	"github.com/opencord/voltha-lib-go/v2/pkg/kafka"
	"github.com/opencord/voltha-protos/v2/go/voltha"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoreProxyImplementsAdapterIfCoreProxy(t *testing.T) {
	proxy := &CoreProxy{}

	if _, ok := interface{}(proxy).(adapterIf.CoreProxy); !ok {
		t.Error("common CoreProxy does not implement adapterif.CoreProxy interface")
	}

}

type InvokeRpcArgs struct {
	rpc string
	toTopic *kafka.Topic
	replyToTopic *kafka.Topic
	waitForResponse bool
	key string
	parentDeviceId string
	kvArgs map[int]interface{}
}


type InvokeRpcSpy struct {
	CallCount int
	Calls     map[int]InvokeRpcArgs
}

type mockKafkaIcProxy struct {
	InvokeRpcSpy
}

func (s *mockKafkaIcProxy) Start() error { return nil }
func (s *mockKafkaIcProxy) DeleteTopic(topic kafka.Topic) error { return nil }
func (s *mockKafkaIcProxy) DeviceDiscovered(deviceId string, deviceType string, parentId string, publisher string) error {return nil}
func (s *mockKafkaIcProxy) Stop() {}
func (s *mockKafkaIcProxy) InvokeRPC(ctx context.Context, rpc string, toTopic *kafka.Topic, replyToTopic *kafka.Topic, waitForResponse bool, key string, kvArgs ...*kafka.KVArg) (bool, *any.Any) {
	s.InvokeRpcSpy.CallCount++

	args := make(map[int]interface{}, 4)
	for k, v := range kvArgs {
		args[k] = v
	}

	s.Calls[s.CallCount] = InvokeRpcArgs{
		rpc:rpc,
		toTopic: toTopic,
		replyToTopic: replyToTopic,
		waitForResponse: waitForResponse,
		key: key,
		kvArgs: args,
	}
	return false, &any.Any{}
}
func (s *mockKafkaIcProxy) SubscribeWithRequestHandlerInterface(topic kafka.Topic, handler interface{}) error { return nil }
func (s *mockKafkaIcProxy) SubscribeWithDefaultRequestHandler(topic kafka.Topic, initialOffset int64) error { return nil }
func (s *mockKafkaIcProxy) UnSubscribeFromRequestHandler(topic kafka.Topic) error { return nil }

func TestCoreProxy_GetChildDevice_sn(t *testing.T) {

	var mockKafkaIcProxy = mockKafkaIcProxy{
		InvokeRpcSpy{
			Calls: make(map[int]InvokeRpcArgs),
		},
	}

	proxy := NewCoreProxy(&mockKafkaIcProxy, "testAdapterTopic", "testCoreTopic")

	kwargs := make(map[string]interface{})
	kwargs["serial_number"] = "TEST00000000001"

	parentDeviceId := "aabbcc"
	proxy.GetChildDevice(context.TODO(), parentDeviceId, kwargs)

	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.CallCount, 1)
	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.Calls[1].rpc, "GetChildDevice")
	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.Calls[1].toTopic, &kafka.Topic{Name:"testCoreTopic"})
	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.Calls[1].replyToTopic, &kafka.Topic{Name:"testAdapterTopic"})
	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.Calls[1].waitForResponse, true)
	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.Calls[1].key, parentDeviceId)
	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.Calls[1].kvArgs[0], &kafka.KVArg{Key:"device_id", Value:&voltha.ID{Id: parentDeviceId}})
	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.Calls[1].kvArgs[1], &kafka.KVArg{Key:"serial_number", Value: &ic.StrType{Val: kwargs["serial_number"].(string)}})
}