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
package stats

import (
	"context"
	"time"
)

type CollectorName string

const (
	OltAdapter CollectorName = "adapter_olt"
	OnuAdapter CollectorName = "adapter_onu"
	VCore      CollectorName = "rw_core"
)

func (s CollectorName) String() string {
	switch s {
	case OltAdapter:
		return "adapter_olt"
	case OnuAdapter:
		return "adapter_onu"
	case VCore:
		return "rw_core"
	}
	return "unknown"
}

type StatsManager interface {
	// Start starts the statistics manager with name and makes the collected stats available at port p.
	Start(ctx context.Context, p int, name CollectorName)

	//CountForDevice counts the number of times the counterName happens for device devId with serial number sn. Each call to Count increments it by one.
	CountForDevice(devId, sn string, counterName DeviceCounter)

	//AddForDevice adds val to counter.
	AddForDevice(devId, sn string, counter DeviceCounter, val float64)

	//CollectDurationForDevice calculates the duration from startTime to time.Now() for device devID with serial number sn.
	CollectDurationForDevice(devID, sn string, dName NonDeviceDuration, startTime time.Time)

	//Count counts the number of times the counterName happens. Each call to Count increments it by one.
	Count(counter NonDeviceCounter)

	//Add adds val to counter.
	Add(counter NonDeviceCounter, val float64)

	//CollectDuration calculates the duration from startTime to time.Now().
	CollectDuration(dName NonDeviceDuration, startTime time.Time)
}
