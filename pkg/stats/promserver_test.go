/*
* Copyright 2018-2024 Open Networking Foundation (ONF) and the ONF Contributors

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
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Check how to reset the prom counters and histogram
func TestPromStatsServer_Start(t *testing.T) {
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()

	testPort := 34201

	StatsServer.Start(serverCtx, testPort, VCore)

	//give time to the prom server to start
	time.Sleep(time.Millisecond * 300)

	StatsServer.Count(NumErrorsWritingToBus)
	StatsServer.Count(NumErrorsWritingToBus)

	StatsServer.CountForDevice("dev1", "serial1", NumOnusActivated)
	StatsServer.CountForDevice("dev1", "serial1", NumOnusActivated)
	StatsServer.CountForDevice("dev1", "serial1", NumOnusActivated)

	StatsServer.Add(NumCoreRpcErrors, 4.0)

	StatsServer.AddForDevice("dev2", "serial2", NumDiscoveriesReceived, 56)

	startTime := time.Now()

	time.Sleep(100 * time.Millisecond)
	StatsServer.CollectDurationForDevice("dev3", "sn3", OnuDiscoveryProcTime, startTime)

	StatsServer.CollectDuration(DBWriteTime, startTime)

	clientCtx, clientCancel := context.WithTimeout(context.Background(), time.Second)
	defer clientCancel()

	req, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:%d/metrics", testPort), nil)
	require.NoError(t, err)
	req = req.WithContext(clientCtx)

	client := http.DefaultClient
	res, err := client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Contains(t, string(bodyBytes), `voltha_rw_core_counters{counter="bus_write_errors_total"} 2`)
	assert.Contains(t, string(bodyBytes), `voltha_rw_core_device_counters{counter="activated_onus_total",device_id="dev1",serial_no="serial1"} 3`)
	assert.Contains(t, string(bodyBytes), `voltha_rw_core_counters{counter="core_rpc_errors_total"} 4`)
	assert.Contains(t, string(bodyBytes), `voltha_rw_core_device_counters{counter="discoveries_received_total",device_id="dev2",serial_no="serial2"} 56`)
	assert.Contains(t, string(bodyBytes), `voltha_rw_core_device_durations_bucket{device_id="dev3",duration="onu_discovery_proc_time",serial_no="sn3",le="300"} 1`)
}
