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

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/opencord/voltha-lib-go/v7/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	cPrefix = "voltha"
)

type PromStatsServer struct {
	// To hold the counters which are specific to devices
	devCounters *prometheus.CounterVec
	// To hold the counters which are NOT tied to specific devices
	otherCounters *prometheus.CounterVec
	// To hold the durations which are specific to devices
	devDurations *prometheus.HistogramVec
	// To hold the durations which are NOT tied to specific to devices
	otherDurations *prometheus.HistogramVec
}

var StatsServer = PromStatsServer{}

// Start starts the statistics manager with name and makes the collected stats available
// at port p. All the statistics collected by this collector will be appended with "voltha_(name)_"
// when they appear in Prometheus. The function starts a prometheus HTTP listener in the background and does not return
// any errors, the listener is stopped on context cancellation
func (ps *PromStatsServer) Start(ctx context.Context, p int, name CollectorName) {
	//log.SetLogger(logging.New())
	ps.initializeCollectors(ctx, name)

	logger.Infow(ctx, "Starting Statistics HTTP Server", log.Fields{"listeningPort": p})

	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: fmt.Sprintf(":%d", p), Handler: nil}

	go func() {
		<-ctx.Done()
		logger.Infof(ctx, "Shutting down the Statistics HTTP server")
		err := server.Shutdown(ctx)
		if err != nil {
			logger.Errorw(ctx, "Statistics server shutting down failure", log.Fields{"error": err})
		}
	}()

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Errorw(ctx, "Starting Statistics HTTP server error", log.Fields{"error": err})
		}
	}()

	logger.Infow(ctx, "Started Prometheus listener for statistics on port", log.Fields{"port": p})
}

func (ps *PromStatsServer) initializeCollectors(ctx context.Context, name CollectorName) {
	logger.Infof(ctx, "Initializing statistics collector")

	collectorName := cPrefix + "_" + name.String()

	var (
		// in milliseconds
		defBuckets = []float64{2, 5, 10, 25, 50, 100, 300, 1000, 5000}
	)

	ps.devCounters = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: collectorName,
			Name:      "device_counters",
			Help:      "Device specific counters",
		},
		[]string{"device_id", "serial_no", "counter"},
	)

	ps.otherCounters = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: collectorName,
			Name:      "counters",
			Help:      "Non device counters",
		},
		[]string{"counter"},
	)

	ps.devDurations = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: collectorName,
			Name:      "device_durations",
			Help:      "Time taken in ms to complete a specific task for a specific device",
			Buckets:   defBuckets,
		},
		[]string{"device_id", "serial_no", "duration"},
	)

	ps.otherDurations = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: collectorName,
			Name:      "durations",
			Help:      "Time taken in ms to complete a specific task not tied to a device",
			Buckets:   defBuckets,
		},
		[]string{"duration"},
	)

	prometheus.MustRegister(ps.devCounters)
	prometheus.MustRegister(ps.otherCounters)
	prometheus.MustRegister(ps.devDurations)
	prometheus.MustRegister(ps.otherDurations)
}

// CountForDevice counts the number of times the counterName happens for device devId with serial number sn. Each call to Count increments it by one.
func (ps *PromStatsServer) CountForDevice(devId, sn string, counterName DeviceCounter) {
	if ps.devCounters != nil {
		ps.devCounters.WithLabelValues(devId, sn, counterName.String()).Inc()
	}
}

// AddForDevice adds val to counter.
func (ps *PromStatsServer) AddForDevice(devId, sn string, counter DeviceCounter, val float64) {
	if ps.devCounters != nil {
		ps.devCounters.WithLabelValues(devId, sn, counter.String()).Add(val)
	}
}

// CollectDurationForDevice calculates the duration from startTime to time.Now() for device devID with serial number sn.
func (ps *PromStatsServer) CollectDurationForDevice(devID, sn string, dName DeviceDuration, startTime time.Time) {
	if ps.otherDurations != nil {
		timeSpent := time.Since(startTime)
		ps.devDurations.WithLabelValues(devID, sn, dName.String()).Observe(float64(timeSpent.Milliseconds()))
	}
}

// Count counts the number of times the counterName happens. Each call to Count increments it by one.
func (ps *PromStatsServer) Count(counter NonDeviceCounter) {
	if ps.otherCounters != nil {
		ps.otherCounters.WithLabelValues(counter.String()).Inc()
	}
}

// Add adds val to counter.
func (ps *PromStatsServer) Add(counter NonDeviceCounter, val float64) {
	if ps.otherCounters != nil {
		ps.otherCounters.WithLabelValues(counter.String()).Add(val)
	}
}

// CollectDuration calculates the duration from startTime to time.Now().
func (ps *PromStatsServer) CollectDuration(dName NonDeviceDuration, startTime time.Time) {
	if ps.otherDurations != nil {
		timeSpent := time.Since(startTime)
		ps.otherDurations.WithLabelValues(dName.String()).Observe(float64(timeSpent.Milliseconds()))
	}
}

// [EOF] - 20231222: Ignore, this triage patch will be abandoned
