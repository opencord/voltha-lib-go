package meters

import (
	"context"
	"github.com/opencord/voltha-protos/v4/go/openflow_13"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMeters_TestTcontType1(t *testing.T) {
	//tcont-type-1
	meterConfig := &openflow_13.OfpMeterConfig{
		MeterId: 1,
		Bands: []*openflow_13.OfpMeterBandHeader{
			{
				Rate:      10000,
				BurstSize: 0,
			},
		},
	}
	shapingInfo := GetTrafficShapingInfo(context.Background(), meterConfig)
	assert.Equal(t, uint32(10000), shapingInfo.Gir)

	//tcont-type-1
	meterConfig = &openflow_13.OfpMeterConfig{
		MeterId: 1,
		Bands: []*openflow_13.OfpMeterBandHeader{
			{
				Rate:      10000,
				BurstSize: 0,
			},
			{
				Rate:      10000,
				BurstSize: 0,
			},
		},
	}
	shapingInfo = GetTrafficShapingInfo(context.Background(), meterConfig)
	assert.Equal(t, uint32(10000), shapingInfo.Pir)
	assert.Equal(t, uint32(10000), shapingInfo.Gir)
}

func TestMeters_TestTcontType2and3(t *testing.T) {
	meterConfig := &openflow_13.OfpMeterConfig{
		MeterId: 1,
		Bands: []*openflow_13.OfpMeterBandHeader{
			{
				Rate:      10000,
				BurstSize: 2000,
			},
			{
				Rate:      30000,
				BurstSize: 3000,
			},
		},
	}
	shapingInfo := GetTrafficShapingInfo(context.Background(), meterConfig)
	assert.Equal(t, uint32(30000), shapingInfo.Pir)
	assert.Equal(t, uint32(10000), shapingInfo.Cir)
}

func TestMeters_TestTcontType4(t *testing.T) {
	meterConfig := &openflow_13.OfpMeterConfig{
		MeterId: 1,
		Bands: []*openflow_13.OfpMeterBandHeader{
			{
				Rate:      10000,
				BurstSize: 1000,
			},
		},
	}
	shapingInfo := GetTrafficShapingInfo(context.Background(), meterConfig)
	assert.Equal(t, uint32(10000), shapingInfo.Pir)
}

func TestMeters_TestTcontType5(t *testing.T) {
	meterConfig := &openflow_13.OfpMeterConfig{
		MeterId: 1,
		Bands: []*openflow_13.OfpMeterBandHeader{
			{
				Rate:      10000,
				BurstSize: 0,
			},
			{
				Rate:      20000,
				BurstSize: 4000,
			},
			{
				Rate:      30000,
				BurstSize: 5000,
			},
		},
	}
	shapingInfo := GetTrafficShapingInfo(context.Background(), meterConfig)
	assert.Equal(t, uint32(30000), shapingInfo.Pir)
	assert.Equal(t, uint32(10000), shapingInfo.Gir)
	assert.Equal(t, uint32(20000), shapingInfo.Cir)
	assert.Equal(t, uint32(5000), shapingInfo.Pbs)
}
