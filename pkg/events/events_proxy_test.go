/*
 * Copyright 2020-present Open Networking Foundation

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

package events

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/opencord/voltha-lib-go/v4/pkg/kafka"
	mock_kafka "github.com/opencord/voltha-lib-go/v4/pkg/mocks/kafka"
	"github.com/opencord/voltha-protos/v4/go/common"
	"github.com/opencord/voltha-protos/v4/go/voltha"
	"github.com/stretchr/testify/assert"
)

func checkSimilarEvent(e1, e2 voltha.Event) bool {
	if reflect.DeepEqual((e1.EventType).(*voltha.Event_RpcEvent), (e2.EventType).(*voltha.Event_RpcEvent)) {
		if (e1.Header.Id != e2.Header.Id) &&
			e1.Header.Category == e2.Header.Category &&
			e1.Header.SubCategory == e2.Header.SubCategory &&
			e1.Header.Type == e2.Header.Type &&
			e1.Header.TypeVersion == e2.Header.TypeVersion &&
			proto.Equal(e1.Header.RaisedTs, e2.Header.RaisedTs) {
			return true
		}
	}
	return false
}

func TestEventProxyReceiveAndSendMessage(t *testing.T) {
	// Init Kafka client
	cTkc := mock_kafka.NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}

	// Init Event Proxy
	ep := NewEventProxy(MsgClient(cTkc), MsgTopic(topic))
	go ep.Start()
	defer ep.Stop()
	testCh := make(chan bool)
	maxWait := 30 * time.Millisecond

	eventMsg := &voltha.RPCEvent{
		Rpc:         "dummy",
		OperationId: "dummy",
		ResourceId:  "dummy",
		Service:     "dummy",
		StackId:     "dummy",
		Status: &common.OperationResp{
			Code: common.OperationResp_OPERATION_FAILURE,
		},
		Description: "dummy",
		Context:     nil,
	}
	var event voltha.Event
	raisedTS := time.Now().Unix()
	event.Header, _ = ep.getEventHeader("RPC_ERROR_RAISE_EVENT", voltha.EventCategory_COMMUNICATION, nil, voltha.EventType_RPC_EVENT, raisedTS)
	event.EventType = &voltha.Event_RpcEvent{RpcEvent: eventMsg}
	timer := time.NewTimer(maxWait)
	defer timer.Stop()
	go func() {
		for {
			select {
			case <-time.After(1 * time.Millisecond):
				if checkSimilarEvent(*ep.eventQueue.current.event, event) {
					testCh <- true
					return
				}
			case <-timer.C:
				testCh <- false
				return
			}
		}
	}()
	err := ep.SendRPCEvent(context.Background(), "RPC_ERROR_RAISE_EVENT", eventMsg, voltha.EventCategory_COMMUNICATION,
		nil, raisedTS)
	assert.Nil(t, err)
	// Check that event is inserted in the queue at last
	ep.eventQueue.mutex.RLock()
	eventPushed := checkSimilarEvent(*ep.eventQueue.last.event, event)
	assert.True(t, eventPushed)
	ep.eventQueue.mutex.RUnlock()
	res := <-testCh
	assert.True(t, res)
}

func TestEventProxyReceiveAndSendMultipleMessagesAtDifferentTime(t *testing.T) {
	// Init Kafka client
	cTkc := mock_kafka.NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}

	// Init Event Proxy
	ep := NewEventProxy(MsgClient(cTkc), MsgTopic(topic))
	go ep.Start()
	defer ep.Stop()
	testCh := make(chan bool)
	maxWait := 30 * time.Millisecond

	eventMsg := &voltha.RPCEvent{
		Rpc:         "dummy",
		OperationId: "dummy",
		ResourceId:  "dummy",
		Service:     "dummy",
		StackId:     "dummy",
		Status: &common.OperationResp{
			Code: common.OperationResp_OPERATION_FAILURE,
		},
		Description: "dummy",
		Context:     nil,
	}
	for index := 1; index <= 3; index++ {
		var event voltha.Event
		raisedTS := time.Now().Unix()
		eventMsg.OperationId = fmt.Sprintf("dummy-%d", index)
		event.Header, _ = ep.getEventHeader("RPC_ERROR_RAISE_EVENT", voltha.EventCategory_COMMUNICATION, nil, voltha.EventType_RPC_EVENT, raisedTS)
		event.EventType = &voltha.Event_RpcEvent{RpcEvent: eventMsg}
		timer := time.NewTimer(maxWait)
		defer timer.Stop()
		go func() {
			for {
				select {
				case <-time.After(1 * time.Millisecond):
					if checkSimilarEvent(*ep.eventQueue.current.event, event) {
						testCh <- true
						return
					}
				case <-timer.C:
					testCh <- false
					return
				}
			}
		}()
		err := ep.SendRPCEvent(context.Background(), "RPC_ERROR_RAISE_EVENT", eventMsg, voltha.EventCategory_COMMUNICATION,
			nil, raisedTS)
		assert.Nil(t, err)
		// Check that event is inserted in the queue at last
		ep.eventQueue.mutex.RLock()
		eventPushed := checkSimilarEvent(*ep.eventQueue.last.event, event)
		assert.True(t, eventPushed)
		ep.eventQueue.mutex.RUnlock()
		res := <-testCh
		assert.True(t, res)
	}

}

func TestEventProxyReceiveAndSendMultipleMessagesAtSameTime(t *testing.T) {
	// Init Kafka client
	cTkc := mock_kafka.NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}

	// Init Event Proxy
	ep := NewEventProxy(MsgClient(cTkc), MsgTopic(topic))
	go ep.Start()
	defer ep.Stop()

	var eventMessageList []voltha.Event
	raisedTS := time.Now().Unix()
	for index := 1; index <= 3; index++ {
		eventMsg := &voltha.RPCEvent{
			Rpc:         "dummy",
			OperationId: fmt.Sprintf("dummy-%d", index),
			ResourceId:  "dummy",
			Service:     "dummy",
			StackId:     "dummy",
			Status: &common.OperationResp{
				Code: common.OperationResp_OPERATION_FAILURE,
			},
			Description: "dummy",
			Context:     nil,
		}
		var event voltha.Event
		event.Header, _ = ep.getEventHeader("RPC_ERROR_RAISE_EVENT", voltha.EventCategory_COMMUNICATION, nil, voltha.EventType_RPC_EVENT, raisedTS)
		event.EventType = &voltha.Event_RpcEvent{RpcEvent: eventMsg}
		eventMessageList = append(eventMessageList, event)
	}

	testCh := make(chan bool)
	maxWait := 30000 * time.Microsecond
	timer := time.NewTimer(maxWait)
	defer timer.Stop()

	go func() {
		for {
			select {
			case <-time.After(1 * time.Microsecond):
				if ep.eventQueue.current != nil {
					//Waiting for the last event to be sent out from the queue
					if checkSimilarEvent(*ep.eventQueue.current.event, eventMessageList[2]) {
						testCh <- true
						return
					}
				}
			case <-timer.C:
				testCh <- false
				return
			}
		}
	}()
	for index := 1; index <= 3; index++ {
		err := ep.SendRPCEvent(context.Background(), "RPC_ERROR_RAISE_EVENT", eventMessageList[index-1].EventType.(*voltha.Event_RpcEvent).RpcEvent, voltha.EventCategory_COMMUNICATION,
			nil, raisedTS)
		assert.Nil(t, err)
		// Check that event is inserted in the queue at last
		ep.eventQueue.mutex.RLock()
		eventPushed := checkSimilarEvent(*ep.eventQueue.last.event, eventMessageList[index-1])
		assert.True(t, eventPushed)
		ep.eventQueue.mutex.RUnlock()
	}
	res := <-testCh
	assert.True(t, res)
}

func TestEventProxyStop(t *testing.T) {
	// Init Kafka client
	cTkc := mock_kafka.NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}

	// Init Event Proxy
	ep := NewEventProxy(MsgClient(cTkc), MsgTopic(topic))
	go ep.Start()

	var eventMessageList []voltha.Event
	raisedTS := time.Now().Unix()
	for index := 1; index <= 3; index++ {
		eventMsg := &voltha.RPCEvent{
			Rpc:         "dummy",
			OperationId: fmt.Sprintf("dummy-%d", index),
			ResourceId:  "dummy",
			Service:     "dummy",
			StackId:     "dummy",
			Status: &common.OperationResp{
				Code: common.OperationResp_OPERATION_FAILURE,
			},
			Description: "dummy",
			Context:     nil,
		}
		var event voltha.Event
		event.Header, _ = ep.getEventHeader("RPC_ERROR_RAISE_EVENT", voltha.EventCategory_COMMUNICATION, nil, voltha.EventType_RPC_EVENT, raisedTS)
		event.EventType = &voltha.Event_RpcEvent{RpcEvent: eventMsg}
		eventMessageList = append(eventMessageList, event)
	}

	for index := 1; index <= 3; index++ {
		err := ep.SendRPCEvent(context.Background(), "RPC_ERROR_RAISE_EVENT", eventMessageList[index-1].EventType.(*voltha.Event_RpcEvent).RpcEvent, voltha.EventCategory_COMMUNICATION,
			nil, raisedTS)
		assert.Nil(t, err)
		// Check that event is inserted in the queue at last
		ep.eventQueue.mutex.RLock()
		eventPushed := checkSimilarEvent(*ep.eventQueue.last.event, eventMessageList[index-1])
		assert.True(t, eventPushed)
		ep.eventQueue.mutex.RUnlock()
	}

	testCh := make(chan bool)
	maxWait := 30000 * time.Microsecond
	timer := time.NewTimer(maxWait)
	defer timer.Stop()

	go func() {
		ep.Stop()
	}()

	go func() {
		for {
			select {
			case <-time.After(1 * time.Microsecond):
				if ep.eventQueue.current != nil {
					//Waiting for the queue to be empty
					if checkSimilarEvent(*ep.eventQueue.current.event, eventMessageList[2]) {
						if ep.stopped && ep.eventQueue.current.prev == nil {
							testCh <- true
							return
						}
					}
				}
			case <-timer.C:
				testCh <- false
				return
			}
		}
	}()

	res := <-testCh
	assert.True(t, res)
}
