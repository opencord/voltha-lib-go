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

//Package adaptercore provides the utility for olt devices, flows and statistics
package common

import (
	"github.com/opencord/voltha-lib-go/v2/pkg/log"
	"github.com/opencord/voltha-protos/v2/go/voltha"
	"strings"
)

const (
	TS          = "ts"
	SLICE       = "slice"
	ADD         = "add"
	REMOVE      = "remove"
	UPDATE      = "update"
	ALL         = "all"
	KPIEVENT    = "kpi_event"
	DEVICEEVENT = "device_event"
	CONFIGEVENT = "config_event"
)

type EventFilter struct {
	Filters map[string]*voltha.EventFilter
}

func NewEventFilter() *EventFilter {
	return &EventFilter{Filters: make(map[string]*voltha.EventFilter)}
}

func (ef *EventFilter) FilterEvent(event *voltha.Event) bool {
	if _, ok := ef.Filters[ALL]; ok {
		return ef.checkAnyEvent(event)
	}
	switch event.Header.Type {
	case voltha.EventType_DEVICE_EVENT:
		return ef.checkDeviceEvent(event)
	case voltha.EventType_KPI_EVENT2:
		return ef.checkKpiEvent(event)
	case voltha.EventType_CONFIG_EVENT:
		return ef.checkConfigEvent(event)
	default:
		log.Errorw("unknown-event-type", log.Fields{"event-type": event.Header.Type})
		return false
	}
}

func (ef *EventFilter) AddRemoveFilters(filter *voltha.EventFilter, remove bool) {
	if !remove {
		ef.Filters[filter.EventType] = filter
		return
	}
	delete(ef.Filters, filter.EventType)
}

func (ef *EventFilter) getEventData(event *voltha.Event) map[voltha.EventFilterRuleKey_EventFilterRuleType]string {
	log.Debugw("Received-event", log.Fields{"event": event})
	eventData := make(map[voltha.EventFilterRuleKey_EventFilterRuleType]string)
	eventData[voltha.EventFilterRuleKey_category] = strings.ToLower(event.Header.Category.String())
	eventData[voltha.EventFilterRuleKey_sub_category] = strings.ToLower(event.Header.SubCategory.String())
	switch event.Header.Type {
	case voltha.EventType_DEVICE_EVENT:
		eventType := strings.Split(event.EventType.(*voltha.Event_DeviceEvent).DeviceEvent.DeviceEventName, "_")
		eventType = eventType[:len(eventType)-2]
		eventData[voltha.EventFilterRuleKey_device_event_type] = strings.Join(eventType, "_")
	case voltha.EventType_KPI_EVENT2:
		kpiEventType := event.EventType.(*voltha.Event_KpiEvent2).KpiEvent2.Type
		if kpiEventType == voltha.KpiEventType_slice {
			eventData[voltha.EventFilterRuleKey_kpi_event_type] = SLICE
		} else if kpiEventType == voltha.KpiEventType_ts {
			eventData[voltha.EventFilterRuleKey_kpi_event_type] = TS
		}
	case voltha.EventType_CONFIG_EVENT:
		configEventType := event.EventType.(*voltha.Event_ConfigEvent).ConfigEvent.Type
		if configEventType == voltha.ConfigEventType_add {
			eventData[voltha.EventFilterRuleKey_config_event_type] = ADD
		} else if configEventType == voltha.ConfigEventType_remove {
			eventData[voltha.EventFilterRuleKey_config_event_type] = REMOVE
		} else if configEventType == voltha.ConfigEventType_update {
			eventData[voltha.EventFilterRuleKey_config_event_type] = UPDATE
		}
	}
	log.Debugw("event-data-created", log.Fields{"event-data": eventData})
	return eventData
}

func (ef *EventFilter) checkAnyEvent(event *voltha.Event) bool {
	if filter, ok := ef.Filters[ALL]; ok {
		log.Info("Fetched filter for event type ANY", log.Fields{"filter": filter})
		if !filter.Enable {
			log.Debugw("Allowing this event as the filter is disabled", log.Fields{"enable": filter.Enable})
			return true
		}
		eventData := ef.getEventData(event)
		if filter.Rules[0].Key == voltha.EventFilterRuleKey_filter_all && filter.Rules[0].Value == "true" {
			log.Debugw("Filter all rule set for all events", log.Fields{"filter-all-rule": filter.Rules[0].Value})
			return false
		}
		for _, rule := range filter.Rules {
			if eventData[rule.Key] != rule.Value {
				log.Debugw("Rules did not match", log.Fields{"filter-rule-value": rule.Value, "event-data": eventData[rule.Key]})
				return true
			}
		}
	} else {
		log.Debugw("No filters presents", log.Fields{"filters": ef.Filters[ALL]})
		return true
	}
	log.Debug("Filtering out this event")
	return false
}

func (ef *EventFilter) checkDeviceEvent(event *voltha.Event) bool {
	if filter, ok := ef.Filters[DEVICEEVENT]; ok {
		log.Debugw("Fetched device event filter", log.Fields{"kpi-event-filter": filter})
		if !filter.Enable {
			log.Debugw("Allowing this event as the filter is disabled", log.Fields{"enable": filter.Enable})
			return true
		}
		eventData := ef.getEventData(event)
		if filter.Rules[0].Key == voltha.EventFilterRuleKey_filter_all && filter.Rules[0].Value == "true" {
			log.Debugw("Filter all rule set for KPI events", log.Fields{"filter-all-rule": filter.Rules[0].Value})
			return false
		}
		for _, rule := range filter.Rules {
			if eventData[rule.Key] != rule.Value {
				log.Debugw("Rules did not match", log.Fields{"filter-rule-value": rule.Value, "event-data": eventData[rule.Key]})
				return true
			}
		}
	} else {
		log.Debugw("No filters present for device event", log.Fields{"filters": ef.Filters[KPIEVENT]})
		return true
	}
	log.Debug("Filtering out this event")
	return false
}

func (ef *EventFilter) checkKpiEvent(event *voltha.Event) bool {
	if filter, ok := ef.Filters[KPIEVENT]; ok {
		log.Debugw("Fetched KPI event filter", log.Fields{"kpi-event-filter": filter})
		if !filter.Enable {
			log.Debugw("Allowing this event as the filter is disabled", log.Fields{"enable": filter.Enable})
			return true
		}
		if filter.Rules[0].Key == voltha.EventFilterRuleKey_filter_all && filter.Rules[0].Value == "true" {
			log.Debugw("Filter all rule set for KPI events", log.Fields{"filter-all-rule": filter.Rules[0].Value})
			return false
		}
		eventData := ef.getEventData(event)
		for _, rule := range filter.Rules {
			if eventData[rule.Key] != rule.Value {
				log.Debugw("Rules did not match", log.Fields{"filter-rule-value": rule.Value, "event-data": eventData[rule.Key]})
				return true
			}
		}
	} else {
		log.Debugw("No filters present for KPI event", log.Fields{"filters": ef.Filters[KPIEVENT]})
		return true
	}
	log.Debug("Filtering out this event")
	return false
}

func (ef *EventFilter) checkConfigEvent(event *voltha.Event) bool {
	if filter, ok := ef.Filters[CONFIGEVENT]; ok {
		log.Debugw("Fetched config event filter", log.Fields{"kpi-event-filter": filter})
		if !filter.Enable {
			log.Debugw("Allowing this event as the filter is disabled", log.Fields{"enable": filter.Enable})
			return true
		}
		if filter.Rules[0].Key == voltha.EventFilterRuleKey_filter_all && filter.Rules[0].Value == "true" {
			log.Debugw("Filter all rule set for config events", log.Fields{"filter-all-rule": filter.Rules[0].Value})
			return false
		}
		eventData := ef.getEventData(event)

		for _, rule := range filter.Rules {
			if eventData[rule.Key] != rule.Value {
				log.Debugw("Rules did not match", log.Fields{"filter-rule-value": rule.Value, "event-data": eventData[rule.Key]})
				return true
			}
		}
	} else {
		log.Debugw("No filters present for config event", log.Fields{"filters": ef.Filters[CONFIGEVENT]})
		return true
	}
	log.Debug("Filtering out this event")
	return false
}
