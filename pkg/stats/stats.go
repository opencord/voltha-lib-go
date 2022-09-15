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

// The unit for all the time based metrics are in milli seconds

const (
	// OLT Device stats
	//-----------------
	// Number of ONU Discovery messages received by the OLT adapter
	NumDiscoveriesReceived DeviceCounter = "NumDiscoveriesReceived"
	// Number of ONUs successfully activated by the OLT adapter
	NumOnusActivated DeviceCounter = "NumOnusActivated"

	// ONT Device stats
	//-----------------
	// Number of times transmission retries for OMCI messages were done for a specific ONU
	OmciCCTxRetries DeviceCounter = "OmciCCTxRetries"
	// Number of times transmission timeouts for OMCI messages happened for a specific ONU
	OmciCCTxTimeouts DeviceCounter = "OmciCCTxTimeouts"

	// Other not device specific stats
	//--------------------------------
	// Number of times writing to the message bus failed, could be collected by adapters as well as vCore
	NumErrorsWritingToBus NonDeviceCounter = "NumErrorsWritingToBus"
	// Number of times rpc calls to the vCore resulted in errors at the adapters
	NumCoreRpcErrors NonDeviceCounter = "NumCoreRpcErrors"
	// Number of times rpc calls to the adapters resulted in errors at the vCore
	NumAdapterRpcErrors NonDeviceCounter = "NumAdapterRpcErrors"

	// OLT Device durations
	//---------------------
	// Time taken at the OLT adapter to process successfully an ONU Discovery message received
	OnuDiscoveryProcTime DeviceDuration = "OntDiscoveryProcTime"
	// Time taken at the OLT adapter to successfully activate an ONU
	OnuDiscToActivateTime DeviceDuration = "OnuDiscToActivateTime"
	// Time taken at the OLT adapter from the time the ONU Discovery was received to the first flow being pushed for the ONU
	OnuDiscToFlowsPushedTime DeviceDuration = "OnuDiscToFlowsPushedTime"

	// ONU Device durations
	//---------------------

	// Other not device specific durations
	//------------------------------------
	// Time taken to write an entry to the database, could be collected by all the three users of the database
	DBWriteTime NonDeviceDuration = "DBWriteTime"
	// Time taken to read an entry from the database, could be collected by all the three users of the database
	DBReadTime NonDeviceDuration = "DBReadTime"
)

func (s DeviceCounter) String() string {
	switch s {
	case NumDiscoveriesReceived:
		return "NumDiscoveriesReceived"
	case NumOnusActivated:
		return "NumOnusActivated"
	case OmciCCTxRetries:
		return "OmciCCTxRetries"
	case OmciCCTxTimeouts:
		return "OmciCCTxTimeouts"
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
	case OnuDiscToActivateTime:
		return "OnuDiscToActivateTime"
	case OnuDiscToFlowsPushedTime:
		return "OnuDiscToFlowsPushedTime"
	}
	return "unknown"
}
