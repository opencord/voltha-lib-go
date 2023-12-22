/*
* Copyright 2018-2023 Open Networking Foundation (ONF) and the ONF Contributors

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
	NumDiscoveriesReceived DeviceCounter = "discoveries_received_total"
	// Number of ONUs successfully activated by the OLT adapter
	NumOnusActivated DeviceCounter = "activated_onus_total"

	// ONT Device stats
	//-----------------
	// Number of times transmission retries for OMCI messages were done for a specific ONU
	OmciCCTxRetries DeviceCounter = "omci_cc_tx_retries_total"
	// Number of times transmission timeouts for OMCI messages happened for a specific ONU
	OmciCCTxTimeouts DeviceCounter = "omci_cc_tx_timeous_total"

	// Other not device specific stats
	//--------------------------------
	// Number of times writing to the message bus failed, could be collected by adapters as well as vCore
	NumErrorsWritingToBus NonDeviceCounter = "bus_write_errors_total"
	// Number of times rpc calls to the vCore resulted in errors at the adapters
	NumCoreRpcErrors NonDeviceCounter = "core_rpc_errors_total"
	// Number of times rpc calls to the adapters resulted in errors at the vCore
	NumAdapterRpcErrors NonDeviceCounter = "adapter_rpc_errors_total"

	// OLT Device durations
	//---------------------
	// Time taken at the OLT adapter to process successfully an ONU Discovery message received
	OnuDiscoveryProcTime DeviceDuration = "onu_discovery_proc_time"
	// Time taken at the OLT adapter to successfully activate an ONU
	OnuDiscToActivateTime DeviceDuration = "onu_discovery_to_activate_time"
	// Time taken at the OLT adapter from the time the ONU Discovery was received to the first flow being pushed for the ONU
	OnuDiscToFlowsPushedTime DeviceDuration = "onu_disc_to_flows_pushed_time"

	// ONU Device durations
	//---------------------

	// Other not device specific durations
	//------------------------------------
	// Time taken to write an entry to the database, could be collected by all the three users of the database
	DBWriteTime NonDeviceDuration = "db_write_time"
	// Time taken to read an entry from the database, could be collected by all the three users of the database
	DBReadTime NonDeviceDuration = "db_read_time"
)

func (s DeviceCounter) String() string {
	switch s {
	case NumDiscoveriesReceived:
		return "discoveries_received_total"
	case NumOnusActivated:
		return "activated_onus_total"
	case OmciCCTxRetries:
		return "omci_cc_tx_retries_total"
	case OmciCCTxTimeouts:
		return "omci_cc_tx_timeous_total"
	}
	return "unknown"
}

func (s NonDeviceCounter) String() string {
	switch s {
	case NumErrorsWritingToBus:
		return "bus_write_errors_total"
	case NumCoreRpcErrors:
		return "core_rpc_errors_total"
	case NumAdapterRpcErrors:
		return "adapter_rpc_errors_total"
	}
	return "unknown"
}

func (s NonDeviceDuration) String() string {
	switch s {
	case DBWriteTime:
		return "db_write_time"
	case DBReadTime:
		return "db_read_time"
	}
	return "unknown"
}

func (s DeviceDuration) String() string {
	switch s {
	case OnuDiscoveryProcTime:
		return "onu_discovery_proc_time"
	case OnuDiscToActivateTime:
		return "onu_discovery_to_activate_time"
	case OnuDiscToFlowsPushedTime:
		return "onu_disc_to_flows_pushed_time"
	}
	return "unknown"
}
// [EOF] - 20231222: Ignore, this triage patch will be abandoned
