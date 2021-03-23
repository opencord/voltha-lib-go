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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/opencord/voltha-lib-go/v4/pkg/events/eventif"
	"github.com/opencord/voltha-lib-go/v4/pkg/kafka"
	"github.com/opencord/voltha-lib-go/v4/pkg/log"
	"github.com/opencord/voltha-protos/v4/go/voltha"
)

type event struct {
	prev, next       *event
	event            *voltha.Event
	notifyOnComplete chan<- struct{}
}

type EventQueue struct {
	mutex sync.Mutex

	last, current  *event
	lastCompleteCh <-chan struct{}
	eventChannel   chan *voltha.Event
}

type EventProxy struct {
	kafkaClient kafka.Client
	eventTopic  kafka.Topic
	eventQueue  *EventQueue // Let there be queue per event proxy, irregardless of underlying messaging mechanism.
}

func NewEventQueue() *EventQueue {
	ch := make(chan struct{})
	close(ch) // assume the "current" event is already complete
	// Start the routine to wait on the queue
	eq := &EventQueue{lastCompleteCh: ch, eventChannel: make(chan (*voltha.Event))}
	return eq
}

func NewEventProxy(opts ...EventProxyOption) *EventProxy {
	var proxy EventProxy
	for _, option := range opts {
		option(&proxy)
	}
	proxy.eventQueue = NewEventQueue()
	return &proxy
}

type EventProxyOption func(*EventProxy)

func MsgClient(client kafka.Client) EventProxyOption {
	return func(args *EventProxy) {
		args.kafkaClient = client
	}
}

func MsgTopic(topic kafka.Topic) EventProxyOption {
	return func(args *EventProxy) {
		args.eventTopic = topic
	}
}

func (ep *EventProxy) formatId(eventName string) string {
	return fmt.Sprintf("Voltha.openolt.%s.%s", eventName, strconv.FormatInt(time.Now().UnixNano(), 10))
}

func (ep *EventProxy) getEventHeader(eventName string,
	category eventif.EventCategory,
	subCategory *eventif.EventSubCategory,
	eventType eventif.EventType,
	raisedTs int64) (*voltha.EventHeader, error) {
	var header voltha.EventHeader
	if strings.Contains(eventName, "_") {
		eventName = strings.Join(strings.Split(eventName, "_")[:len(strings.Split(eventName, "_"))-2], "_")
	} else {
		eventName = "UNKNOWN_EVENT"
	}
	/* Populating event header */
	header.Id = ep.formatId(eventName)
	header.Category = category
	if subCategory != nil {
		header.SubCategory = *subCategory
	} else {
		header.SubCategory = voltha.EventSubCategory_NONE
	}
	header.Type = eventType
	header.TypeVersion = eventif.EventTypeVersion

	// raisedTs is in nanoseconds
	timestamp, err := ptypes.TimestampProto(time.Unix(0, raisedTs))
	if err != nil {
		return nil, err
	}
	header.RaisedTs = timestamp

	timestamp, err = ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}
	header.ReportedTs = timestamp

	return &header, nil
}

/* Send out rpc events*/
func (ep *EventProxy) SendRPCEvent(ctx context.Context, id string, rpcEvent *voltha.RPCEvent, category eventif.EventCategory, subCategory *eventif.EventSubCategory, raisedTs int64) error {
	if rpcEvent == nil {
		logger.Error(ctx, "Received empty rpc event")
		return errors.New("rpc event nil")
	}
	var event voltha.Event
	var err error
	if event.Header, err = ep.getEventHeader(id, category, subCategory, voltha.EventType_RPC_EVENT, raisedTs); err != nil {
		return err
	}
	event.EventType = &voltha.Event_RpcEvent{RpcEvent: rpcEvent}
	ep.eventQueue.SendEvent(ctx, &event)
	return nil

}

/* Send out device events*/
func (ep *EventProxy) SendDeviceEvent(ctx context.Context, deviceEvent *voltha.DeviceEvent, category eventif.EventCategory, subCategory eventif.EventSubCategory, raisedTs int64) error {
	if deviceEvent == nil {
		logger.Error(ctx, "Recieved empty device event")
		return errors.New("Device event nil")
	}
	var event voltha.Event
	var de voltha.Event_DeviceEvent
	var err error
	de.DeviceEvent = deviceEvent
	if event.Header, err = ep.getEventHeader(deviceEvent.DeviceEventName, category, &subCategory, voltha.EventType_DEVICE_EVENT, raisedTs); err != nil {
		return err
	}
	event.EventType = &de
	if err := ep.sendEvent(ctx, &event); err != nil {
		logger.Errorw(ctx, "Failed to send device event to KAFKA bus", log.Fields{"device-event": deviceEvent})
		return err
	}
	logger.Infow(ctx, "Successfully sent device event KAFKA", log.Fields{"Id": event.Header.Id, "Category": event.Header.Category,
		"SubCategory": event.Header.SubCategory, "Type": event.Header.Type, "TypeVersion": event.Header.TypeVersion,
		"ReportedTs": event.Header.ReportedTs, "ResourceId": deviceEvent.ResourceId, "Context": deviceEvent.Context,
		"DeviceEventName": deviceEvent.DeviceEventName})

	return nil

}

// SendKpiEvent is to send kpi events to voltha.event topic
func (ep *EventProxy) SendKpiEvent(ctx context.Context, id string, kpiEvent *voltha.KpiEvent2, category eventif.EventCategory, subCategory eventif.EventSubCategory, raisedTs int64) error {
	if kpiEvent == nil {
		logger.Error(ctx, "Recieved empty kpi event")
		return errors.New("KPI event nil")
	}
	var event voltha.Event
	var de voltha.Event_KpiEvent2
	var err error
	de.KpiEvent2 = kpiEvent
	if event.Header, err = ep.getEventHeader(id, category, &subCategory, voltha.EventType_KPI_EVENT2, raisedTs); err != nil {
		return err
	}
	event.EventType = &de
	if err := ep.sendEvent(ctx, &event); err != nil {
		logger.Errorw(ctx, "Failed to send kpi event to KAFKA bus", log.Fields{"device-event": kpiEvent})
		return err
	}
	logger.Infow(ctx, "Successfully sent kpi event to KAFKA", log.Fields{"Id": event.Header.Id, "Category": event.Header.Category,
		"SubCategory": event.Header.SubCategory, "Type": event.Header.Type, "TypeVersion": event.Header.TypeVersion,
		"ReportedTs": event.Header.ReportedTs, "KpiEventName": "STATS_EVENT"})

	return nil

}

func (ep *EventProxy) sendEvent(ctx context.Context, event *voltha.Event) error {
	logger.Debugw(ctx, "Send event to kafka", log.Fields{"event": event})
	if err := ep.kafkaClient.Send(ctx, event, &ep.eventTopic); err != nil {
		return err
	}
	logger.Debugw(ctx, "Sent event to kafka", log.Fields{"event": event})

	return nil
}

func (ep *EventProxy) EnableLivenessChannel(ctx context.Context, enable bool) chan bool {
	return ep.kafkaClient.EnableLivenessChannel(ctx, enable)
}

func (ep *EventProxy) SendLiveness(ctx context.Context) error {
	return ep.kafkaClient.SendLiveness(ctx)
}

func (ep *EventProxy) ReceiveAndPublishEvents() {
	eq := ep.eventQueue
	for {
		select {
		case event := <-eq.eventChannel:
			eq.mutex.Lock()
			ctx := context.Background()
			if err := ep.sendEvent(ctx, event); err != nil {
				logger.Errorw(ctx, "Failed to send event to KAFKA bus", log.Fields{"event": event})
			} else {
				logger.Debugw(ctx, "Successfully sent RPC event to KAFKA bus", log.Fields{"Id": event.Header.Id, "Category": event.Header.Category,
					"SubCategory": event.Header.SubCategory, "Type": event.Header.Type, "TypeVersion": event.Header.TypeVersion,
					"ReportedTs": event.Header.ReportedTs, "event-type": event.EventType})
			}
			// Notify the next waiting event.
			close(eq.current.notifyOnComplete)
			if eq.current.next != nil {
				eq.current.next.prev = nil
			}
			eq.mutex.Unlock()
		default:

		}
	}
}

func (eq *EventQueue) SendEvent(ctx context.Context, ev *voltha.Event) {
	// add event to the end of the queue and then wait for the turn
	eq.mutex.Lock()
	waitingOn := eq.lastCompleteCh

	ch := make(chan struct{})
	eq.lastCompleteCh = ch
	r := &event{notifyOnComplete: ch, event: ev}

	if eq.last != nil {
		eq.last.next, r.prev = r, eq.last
	}
	eq.last = r
	eq.mutex.Unlock()
	go func() {
		<-waitingOn
		// Previous event has signaled that it is complete.
		// This event now can proceed as the active
		// event
		eq.mutex.Lock()

		eq.current = r
		currentEvent := r.event
		eq.mutex.Unlock()
		// Send current event
		eq.eventChannel <- currentEvent

	}()
}
