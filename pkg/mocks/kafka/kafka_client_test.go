/*
 * Copyright 2019-2024 Open Networking Foundation (ONF) and the ONF Contributors
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

package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/opencord/voltha-lib-go/v7/pkg/kafka"
	"github.com/opencord/voltha-protos/v5/go/core_adapter"
	"github.com/stretchr/testify/assert"
)

func TestKafkaClientCreateTopic(t *testing.T) {
	ctx := context.Background()
	cTkc := NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}
	err := cTkc.CreateTopic(ctx, &topic, 1, 1)
	assert.Nil(t, err)
	err = cTkc.CreateTopic(ctx, &topic, 1, 1)
	assert.NotNil(t, err)
}

func TestKafkaClientDeleteTopic(t *testing.T) {
	cTkc := NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}
	err := cTkc.DeleteTopic(context.Background(), &topic)
	assert.Nil(t, err)
}

func TestKafkaClientSubscribeSend(t *testing.T) {
	cTkc := NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}
	ch, err := cTkc.Subscribe(context.Background(), &topic)
	assert.Nil(t, err)
	assert.NotNil(t, ch)
	testCh := make(chan bool)
	maxWait := 5 * time.Millisecond
	msg := &core_adapter.DeviceReason{
		DeviceId: "1234",
		Reason:   "mock",
	}
	timer := time.NewTimer(maxWait)
	defer timer.Stop()
	go func() {
		select {
		case val, ok := <-ch:
			assert.True(t, ok)
			assert.Equal(t, val, msg)
			testCh <- true
		case <-timer.C:
			testCh <- false
		}
	}()
	err = cTkc.Send(context.Background(), msg, &topic)
	assert.Nil(t, err)
	res := <-testCh
	assert.True(t, res)
}

func TestKafkaClientUnSubscribe(t *testing.T) {
	cTkc := NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}
	ch, err := cTkc.Subscribe(context.Background(), &topic)
	assert.Nil(t, err)
	assert.NotNil(t, ch)
	err = cTkc.UnSubscribe(context.Background(), &topic, ch)
	assert.Nil(t, err)
}

func TestKafkaClientStop(t *testing.T) {
	cTkc := NewKafkaClient()
	topic := kafka.Topic{Name: "myTopic"}
	ch, err := cTkc.Subscribe(context.Background(), &topic)
	assert.Nil(t, err)
	assert.NotNil(t, ch)
	err = cTkc.UnSubscribe(context.Background(), &topic, ch)
	assert.Nil(t, err)
	cTkc.Stop(context.Background())
}
