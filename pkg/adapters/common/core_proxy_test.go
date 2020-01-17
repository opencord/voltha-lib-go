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
	adapterIf "github.com/opencord/voltha-lib-go/v3/pkg/adapters/adapterif"
	"github.com/opencord/voltha-lib-go/v3/pkg/kafka"
	"github.com/opencord/voltha-lib-go/v3/pkg/mocks"
	ic "github.com/opencord/voltha-protos/v3/go/inter_container"
	"github.com/opencord/voltha-protos/v3/go/voltha"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoreProxyImplementsAdapterIfCoreProxy(t *testing.T) {
	proxy := &CoreProxy{}

	if _, ok := interface{}(proxy).(adapterIf.CoreProxy); !ok {
		t.Error("common CoreProxy does not implement adapterif.CoreProxy interface")
	}

}

func TestCoreProxy_GetChildDevice_sn(t *testing.T) {

	var mockKafkaIcProxy = mocks.MockKafkaIcProxy{
		InvokeRpcSpy: mocks.InvokeRpcSpy{
			Calls: make(map[int]mocks.InvokeRpcArgs),
		},
	}

	proxy := NewCoreProxy(&mockKafkaIcProxy, "testAdapterTopic", "testCoreTopic")

	kwargs := make(map[string]interface{})
	kwargs["serial_number"] = "TEST00000000001"

	parentDeviceId := "aabbcc"
	proxy.GetChildDevice(context.TODO(), parentDeviceId, kwargs)

	assert.Equal(t, mockKafkaIcProxy.InvokeRpcSpy.CallCount, 1)
	call := mockKafkaIcProxy.InvokeRpcSpy.Calls[1]
	assert.Equal(t, call.Rpc, "GetChildDevice")
	assert.Equal(t, call.ToTopic, &kafka.Topic{Name: "testCoreTopic"})
	assert.Equal(t, call.ReplyToTopic, &kafka.Topic{Name: "testAdapterTopic"})
	assert.Equal(t, call.WaitForResponse, true)
	assert.Equal(t, call.Key, parentDeviceId)
	assert.Equal(t, call.KvArgs[0], &kafka.KVArg{Key: "device_id", Value: &voltha.ID{Id: parentDeviceId}})
	assert.Equal(t, call.KvArgs[1], &kafka.KVArg{Key: "serial_number", Value: &ic.StrType{Val: kwargs["serial_number"].(string)}})
}
