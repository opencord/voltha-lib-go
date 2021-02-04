/*
 * Copyright 2021-present Open Networking Foundation

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
package onu

import (
	"context"
	"errors"
	"github.com/opencord/voltha-lib-go/v4/pkg/events/eventif"
	"github.com/opencord/voltha-protos/v4/go/voltha"
)

type onuEvents map[onuDevice]OnuDeviceEvent
type onuDevice struct {
	classID uint16
	alarmno uint8
}
type OnuDeviceEvent struct {
	EventName        string
	EventCategory    eventif.EventCategory
	EventSubCategory eventif.EventSubCategory
	EventDescription string
}

const (
	CircuitPackClassID = uint16(6)
	PhysicalPathTerminationPointEthernetUniClassID = uint16(11)
	OnuGClassID = uint16(256)
	AniGClassID = uint16(263)

)
func getOnuEventDetailsByClassIDAndAlarmNo(classID uint16, alarmNo uint8) OnuDeviceEvent {
	onuEventList := func() onuEvents {
		onuEventsList := make(map[onuDevice]OnuDeviceEvent)
		onuEventsList[onuDevice{classID: CircuitPackClassID, alarmno: 0}] = OnuDeviceEvent{EventName: "ONU_EQUIPMENT",
			EventCategory: voltha.EventCategory_EQUIPMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Equipment alarm"}
		onuEventsList[onuDevice{classID: CircuitPackClassID, alarmno: 2}] = OnuDeviceEvent{EventName: "ONU_SELF_TEST_FAIL",
			EventCategory: voltha.EventCategory_EQUIPMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Self-test failure"}
		onuEventsList[onuDevice{classID: CircuitPackClassID, alarmno: 3}] = OnuDeviceEvent{EventName: "ONU_LASER_EOL",
			EventCategory: voltha.EventCategory_EQUIPMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Laser end of life"}
		onuEventsList[onuDevice{classID: CircuitPackClassID, alarmno: 4}] = OnuDeviceEvent{EventName: "ONU_TEMP_YELLOW",
			EventCategory: voltha.EventCategory_ENVIRONMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Temperature yellow"}
		onuEventsList[onuDevice{classID: CircuitPackClassID, alarmno: 5}] = OnuDeviceEvent{EventName: "ONU_TEMP_RED",
			EventCategory: voltha.EventCategory_ENVIRONMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Temperature red"}
		onuEventsList[onuDevice{classID: PhysicalPathTerminationPointEthernetUniClassID, alarmno: 0}] =
			OnuDeviceEvent{EventName: "ONU_Ethernet_UNI", EventCategory: voltha.EventCategory_EQUIPMENT,
				EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "LAN Loss Of Signal"}
		onuEventsList[onuDevice{classID: OnuGClassID, alarmno: 0}] = OnuDeviceEvent{EventName: "ONU_EQUIPMENT",
			EventCategory: voltha.EventCategory_EQUIPMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Equipment alarm"}
		onuEventsList[onuDevice{classID: OnuGClassID, alarmno: 6}] = OnuDeviceEvent{EventName: "ONU_SELF_TEST_FAIL",
			EventCategory: voltha.EventCategory_EQUIPMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Self-test failure"}
		onuEventsList[onuDevice{classID: OnuGClassID, alarmno: 7}] = OnuDeviceEvent{EventName: "ONU_DYING_GASP",
			EventCategory: voltha.EventCategory_EQUIPMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Dying gasp"}
		onuEventsList[onuDevice{classID: OnuGClassID, alarmno: 8}] = OnuDeviceEvent{EventName: "ONU_TEMP_YELLOW",
			EventCategory: voltha.EventCategory_ENVIRONMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Temperature yellow"}
		onuEventsList[onuDevice{classID: OnuGClassID, alarmno: 9}] = OnuDeviceEvent{EventName: "ONU_TEMP_RED",
			EventCategory: voltha.EventCategory_ENVIRONMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Temperature red"}
		onuEventsList[onuDevice{classID: OnuGClassID, alarmno: 10}] = OnuDeviceEvent{EventName: "ONU_VOLTAGE_YELLOW",
			EventCategory: voltha.EventCategory_ENVIRONMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Voltage yellow"}
		onuEventsList[onuDevice{classID: OnuGClassID, alarmno: 11}] = OnuDeviceEvent{EventName: "ONU_VOLTAGE_RED",
			EventCategory: voltha.EventCategory_ENVIRONMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Voltage red"}
		onuEventsList[onuDevice{classID: AniGClassID, alarmno: 0}] = OnuDeviceEvent{EventName: "ONU_LOW_RX_OPTICAL",
			EventCategory: voltha.EventCategory_COMMUNICATION, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Low received optical power"}
		onuEventsList[onuDevice{classID: AniGClassID, alarmno: 1}] = OnuDeviceEvent{EventName: "ONU_HIGH_RX_OPTICAL",
			EventCategory: voltha.EventCategory_COMMUNICATION, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "High received optical power"}
		onuEventsList[onuDevice{classID: AniGClassID, alarmno: 4}] = OnuDeviceEvent{EventName: "ONU_LOW_TX_OPTICAL",
			EventCategory: voltha.EventCategory_COMMUNICATION, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Low transmit optical power"}
		onuEventsList[onuDevice{classID: AniGClassID, alarmno: 5}] = OnuDeviceEvent{EventName: "ONU_HIGH_TX_OPTICAL",
			EventCategory: voltha.EventCategory_COMMUNICATION, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "High transmit optical power"}
		onuEventsList[onuDevice{classID: AniGClassID, alarmno: 6}] = OnuDeviceEvent{EventName: "ONU_LASER_BIAS_CURRENT",
			EventCategory: voltha.EventCategory_EQUIPMENT, EventSubCategory: voltha.EventSubCategory_ONU, EventDescription: "Laser bias current"}
		return onuEventsList
	}()
	return onuEventList[onuDevice{classID: classID, alarmno: alarmNo}]
}

// GetDeviceEventData returns the event data for a device
func GetDeviceEventData(ctx context.Context, classID uint16, alarmNo uint8) (OnuDeviceEvent, error) {
	onuEventDetails := getOnuEventDetailsByClassIDAndAlarmNo(classID, alarmNo)
	if onuEventDetails == (OnuDeviceEvent{}) {
		return onuEventDetails, errors.New("onu Event Detail not found")
	}
	return onuEventDetails, nil
}
