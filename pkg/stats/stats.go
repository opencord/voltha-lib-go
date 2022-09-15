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

type DeviceCounter string
type NonDeviceCounter string

type NonDeviceDuration string
type DeviceDuration string

const (
	// OLT Device stats
	NumDiscoveriesReceived DeviceCounter = "NumDiscoveriesReceived"
	NumOnusActivated       DeviceCounter = "NumOnusActivated"

	// ONT Device stats

	// Other not device specific stats
	NumErrorsWritingToBus NonDeviceCounter = "NumErrorsWritingToBus"
	NumCoreRpcErrors      NonDeviceCounter = "NumCoreRpcErrors"
	NumAdapterRpcErrors   NonDeviceCounter = "NumAdapterRpcErrors"

	// OLT Device durations
	OnuDiscoveryProcTime DeviceDuration = "OntDiscoveryProcTime"

	// Other not device specific durations
	DBWriteTime NonDeviceDuration = "DBWriteTime"
	DBReadTime  NonDeviceDuration = "DBReadTime"
)

func (s DeviceCounter) String() string {
	switch s {
	case NumDiscoveriesReceived:
		return "NumDiscoveriesReceived"
	case NumOnusActivated:
		return "NumOnusActivated"
	}
	return "unknown"
}

func (s NonDeviceCounter) String() string {
	switch s {
	case NumErrorsWritingToBus:
		return "NumErrorsWritingToBus"
	case NumCoreRpcErrors:
		return "NumCoreRpcErrors"
	case NumAdapterRpcErrors:
		return "NumAdapterRpcErrors"
	}
	return "unknown"
}

func (s NonDeviceDuration) String() string {
	switch s {
	case DBWriteTime:
		return "DBWriteTime"
	case DBReadTime:
		return "DBReadTime"
	}
	return "unknown"
}

func (s DeviceDuration) String() string {
	switch s {
	case OnuDiscoveryProcTime:
		return "OnuDiscoveryProcTime"
	}
	return "unknown"
}
