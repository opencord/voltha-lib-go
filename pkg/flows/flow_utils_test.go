/*
 * Copyright 2019-present Open Networking Foundation
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
package flows

import (
	"bytes"
	"context"
	"strings"
	"testing"

	ofp "github.com/opencord/voltha-protos/v3/go/openflow_13"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	timeoutError     error
	taskFailureError error
)

func init() {
	timeoutError = status.Errorf(codes.Aborted, "timeout")
	taskFailureError = status.Error(codes.Internal, "test failure task")
	timeoutError = status.Errorf(codes.Aborted, "timeout")
}

func TestFlowsAndGroups_AddFlow(t *testing.T) {
	ctx := context.Background()
	fg := NewFlowsAndGroups(ctx)
	allFlows := fg.ListFlows(ctx)
	assert.Equal(t, 0, len(allFlows))
	fg.AddFlow(ctx, nil)
	allFlows = fg.ListFlows(ctx)
	assert.Equal(t, 0, len(allFlows))

	fa := &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 1),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|1),
			TunnelId(ctx, uint64(1)),
			EthType(ctx, 0x0800),
			Ipv4Dst(ctx, 0xffffffff),
			IpProto(ctx, 17),
			UdpSrc(ctx, 68),
			UdpDst(ctx, 67),
		},
		Actions: []*ofp.OfpAction{
			PushVlan(ctx, 0x8100),
			SetField(ctx, VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|4000)),
			Output(ctx, uint32(ofp.OfpPortNo_OFPP_CONTROLLER)),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	fg.AddFlow(ctx, flow)

	allFlows = fg.ListFlows(ctx)
	assert.Equal(t, 1, len(allFlows))
	assert.True(t, FlowMatch(ctx, flow, allFlows[0]))
}

func TestFlowsAndGroups_AddGroup(t *testing.T) {
	ctx := context.Background()
	var ga *GroupArgs

	fg := NewFlowsAndGroups(ctx)
	allGroups := fg.ListGroups(ctx)
	assert.Equal(t, 0, len(allGroups))
	fg.AddGroup(ctx, nil)
	allGroups = fg.ListGroups(ctx)
	assert.Equal(t, 0, len(allGroups))

	ga = &GroupArgs{
		GroupId: 10,
		Buckets: []*ofp.OfpBucket{
			{Actions: []*ofp.OfpAction{
				PopVlan(ctx),
				Output(ctx, 1),
			},
			},
		},
	}
	group := MkGroupStat(ctx, ga)
	fg.AddGroup(ctx, group)

	allGroups = fg.ListGroups(ctx)
	assert.Equal(t, 1, len(allGroups))
	assert.Equal(t, ga.GroupId, allGroups[0].Desc.GroupId)
}

func TestFlowsAndGroups_Copy(t *testing.T) {
	ctx := context.Background()
	fg := NewFlowsAndGroups(ctx)
	var fa *FlowArgs
	var ga *GroupArgs

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
		},
		Actions: []*ofp.OfpAction{
			SetField(ctx, VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|10)),
			Output(ctx, 1),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	fg.AddFlow(ctx, flow)

	ga = &GroupArgs{
		GroupId: 10,
		Buckets: []*ofp.OfpBucket{
			{Actions: []*ofp.OfpAction{
				PopVlan(ctx),
				Output(ctx, 1),
			},
			},
		},
	}
	group := MkGroupStat(ctx, ga)
	fg.AddGroup(ctx, group)

	fgCopy := fg.Copy(ctx)

	allFlows := fgCopy.ListFlows(ctx)
	assert.Equal(t, 1, len(allFlows))
	assert.True(t, FlowMatch(ctx, flow, allFlows[0]))

	allGroups := fgCopy.ListGroups(ctx)
	assert.Equal(t, 1, len(allGroups))
	assert.Equal(t, ga.GroupId, allGroups[0].Desc.GroupId)

	fg = NewFlowsAndGroups(ctx)
	fgCopy = fg.Copy(ctx)
	allFlows = fgCopy.ListFlows(ctx)
	allGroups = fgCopy.ListGroups(ctx)
	assert.Equal(t, 0, len(allFlows))
	assert.Equal(t, 0, len(allGroups))
}

func TestFlowsAndGroups_GetFlow(t *testing.T) {
	ctx := context.Background()
	fg := NewFlowsAndGroups(ctx)
	var fa1 *FlowArgs
	var fa2 *FlowArgs
	var ga *GroupArgs

	fa1 = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			Metadata_ofp(ctx, (1000<<32)|1),
			VlanPcp(ctx, 0),
		},
		Actions: []*ofp.OfpAction{
			PopVlan(ctx),
		},
	}
	flow1, err := MkFlowStat(ctx, fa1)
	assert.Nil(t, err)

	fa2 = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 1500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 5),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
		},
		Actions: []*ofp.OfpAction{
			PushVlan(ctx, 0x8100),
			SetField(ctx, VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|1000)),
			SetField(ctx, VlanPcp(ctx, 0)),
			Output(ctx, 2),
		},
	}
	flow2, err := MkFlowStat(ctx, fa2)
	assert.Nil(t, err)

	fg.AddFlow(ctx, flow1)
	fg.AddFlow(ctx, flow2)

	ga = &GroupArgs{
		GroupId: 10,
		Buckets: []*ofp.OfpBucket{
			{Actions: []*ofp.OfpAction{
				PopVlan(ctx),
				Output(ctx, 1),
			},
			},
		},
	}
	group := MkGroupStat(ctx, ga)
	fg.AddGroup(ctx, group)

	gf1 := fg.GetFlow(ctx, 0)
	assert.True(t, FlowMatch(ctx, flow1, gf1))

	gf2 := fg.GetFlow(ctx, 1)
	assert.True(t, FlowMatch(ctx, flow2, gf2))

	gf3 := fg.GetFlow(ctx, 2)
	assert.Nil(t, gf3)

	allFlows := fg.ListFlows(ctx)
	assert.True(t, FlowMatch(ctx, flow1, allFlows[0]))
	assert.True(t, FlowMatch(ctx, flow2, allFlows[1]))
}

func TestFlowsAndGroups_String(t *testing.T) {
	ctx := context.Background()
	fg := NewFlowsAndGroups(ctx)
	var fa *FlowArgs
	var ga *GroupArgs

	str := fg.String(ctx)
	assert.True(t, str == "")

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			Group(ctx, 10),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	fg.AddFlow(ctx, flow)

	ga = &GroupArgs{
		GroupId: 10,
		Buckets: []*ofp.OfpBucket{
			{Actions: []*ofp.OfpAction{
				PopVlan(ctx),
				Output(ctx, 1),
			},
			},
		},
	}
	group := MkGroupStat(ctx, ga)
	fg.AddGroup(ctx, group)

	str = fg.String(ctx)
	assert.True(t, strings.Contains(str, "id: 11819684229970388353"))
	assert.True(t, strings.Contains(str, "group_id: 10"))
	assert.True(t, strings.Contains(str, "oxm_class: OFPXMC_OPENFLOW_BASICOFPXMC_OPENFLOW_BASIC"))
	assert.True(t, strings.Contains(str, "type: OFPXMT_OFB_VLAN_VIDOFPXMT_OFB_VLAN_VID"))
	assert.True(t, strings.Contains(str, "vlan_vid: 4096"))
	assert.True(t, strings.Contains(str, "buckets:"))
}

func TestFlowsAndGroups_AddFrom(t *testing.T) {
	ctx := context.Background()
	fg := NewFlowsAndGroups(ctx)
	var fa *FlowArgs
	var ga *GroupArgs

	str := fg.String(ctx)
	assert.True(t, str == "")

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			Metadata_ofp(ctx, 1000),
			TunnelId(ctx, uint64(1)),
			VlanPcp(ctx, 0),
		},
		Actions: []*ofp.OfpAction{
			PopVlan(ctx),
			Output(ctx, 1),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	fg.AddFlow(ctx, flow)

	ga = &GroupArgs{
		GroupId: 10,
		Buckets: []*ofp.OfpBucket{
			{Actions: []*ofp.OfpAction{
				PopVlan(ctx),
				Output(ctx, 1),
			},
			},
		},
	}
	group := MkGroupStat(ctx, ga)
	fg.AddGroup(ctx, group)

	fg1 := NewFlowsAndGroups(ctx)
	fg1.AddFrom(ctx, fg)

	allFlows := fg1.ListFlows(ctx)
	allGroups := fg1.ListGroups(ctx)
	assert.Equal(t, 1, len(allFlows))
	assert.Equal(t, 1, len(allGroups))
	assert.True(t, FlowMatch(ctx, flow, allFlows[0]))
	assert.Equal(t, group.Desc.GroupId, allGroups[0].Desc.GroupId)
}

func TestDeviceRules_AddFlow(t *testing.T) {
	ctx := context.Background()
	dr := NewDeviceRules(ctx)
	rules := dr.GetRules(ctx)
	assert.True(t, len(rules) == 0)

	dr.AddFlow(ctx, "123456", nil)
	rules = dr.GetRules(ctx)
	assert.True(t, len(rules) == 1)
	val, ok := rules["123456"]
	assert.True(t, ok)
	assert.Equal(t, 0, len(val.ListFlows(ctx)))
	assert.Equal(t, 0, len(val.ListGroups(ctx)))

	fa := &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			Metadata_ofp(ctx, 1000),
			TunnelId(ctx, uint64(1)),
			VlanPcp(ctx, 0),
		},
		Actions: []*ofp.OfpAction{
			PopVlan(ctx),
			Output(ctx, 1),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	dr.AddFlow(ctx, "123456", flow)
	rules = dr.GetRules(ctx)
	assert.True(t, len(rules) == 1)
	val, ok = rules["123456"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(val.ListFlows(ctx)))
	assert.True(t, FlowMatch(ctx, flow, val.ListFlows(ctx)[0]))
	assert.Equal(t, 0, len(val.ListGroups(ctx)))
}

func TestDeviceRules_AddFlowsAndGroup(t *testing.T) {
	ctx := context.Background()
	fg := NewFlowsAndGroups(ctx)
	var fa *FlowArgs
	var ga *GroupArgs

	str := fg.String(ctx)
	assert.True(t, str == "")

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 2000},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			Metadata_ofp(ctx, 1000),
			TunnelId(ctx, uint64(1)),
			VlanPcp(ctx, 0),
		},
		Actions: []*ofp.OfpAction{
			PopVlan(ctx),
			Output(ctx, 1),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	fg.AddFlow(ctx, flow)

	ga = &GroupArgs{
		GroupId: 10,
		Buckets: []*ofp.OfpBucket{
			{Actions: []*ofp.OfpAction{
				PopVlan(ctx),
				Output(ctx, 1),
			},
			},
		},
	}
	group := MkGroupStat(ctx, ga)
	fg.AddGroup(ctx, group)

	dr := NewDeviceRules(ctx)
	dr.AddFlowsAndGroup(ctx, "123456", fg)
	rules := dr.GetRules(ctx)
	assert.True(t, len(rules) == 1)
	val, ok := rules["123456"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(val.ListFlows(ctx)))
	assert.Equal(t, 1, len(val.ListGroups(ctx)))
	assert.True(t, FlowMatch(ctx, flow, val.ListFlows(ctx)[0]))
	assert.Equal(t, 10, int(val.ListGroups(ctx)[0].Desc.GroupId))
}

func TestFlowHasOutPort(t *testing.T) {
	ctx := context.Background()
	var flow *ofp.OfpFlowStats
	assert.False(t, FlowHasOutPort(ctx, flow, 1))

	fa := &FlowArgs{
		KV: OfpFlowModArgs{"priority": 2000},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			Metadata_ofp(ctx, 1000),
			TunnelId(ctx, uint64(1)),
			VlanPcp(ctx, 0),
		},
		Actions: []*ofp.OfpAction{
			PopVlan(ctx),
			Output(ctx, 1),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.True(t, FlowHasOutPort(ctx, flow, 1))
	assert.False(t, FlowHasOutPort(ctx, flow, 2))

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 2000},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
		},
	}
	flow, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowHasOutPort(ctx, flow, 1))
}

func TestFlowHasOutGroup(t *testing.T) {
	ctx := context.Background()
	var flow *ofp.OfpFlowStats
	assert.False(t, FlowHasOutGroup(ctx, flow, 10))

	fa := &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			Group(ctx, 10),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.True(t, FlowHasOutGroup(ctx, flow, 10))
	assert.False(t, FlowHasOutGroup(ctx, flow, 11))

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			Output(ctx, 1),
		},
	}
	flow, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowHasOutGroup(ctx, flow, 1))
}

func TestMatchFlow(t *testing.T) {
	ctx := context.Background()
	assert.False(t, FlowMatch(ctx, nil, nil))
	fa := &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1, "cookie": 38268468, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			Group(ctx, 10),
		},
	}
	flow1, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowMatch(ctx, flow1, nil))

	// different table_id, cookie, flags
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			Group(ctx, 10),
		},
	}
	flow2, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowMatch(ctx, flow1, flow2))
	assert.False(t, FlowMatch(ctx, nil, flow2))

	// no difference
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1, "cookie": 38268468, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.True(t, FlowMatch(ctx, flow1, flow2))

	// different priority
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 501, "table_id": 1, "cookie": 38268468, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			Group(ctx, 10),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowMatch(ctx, flow1, flow2))

	// different table id
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 2, "cookie": 38268468, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowMatch(ctx, flow1, flow2))

	// different cookie
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1, "cookie": 38268467, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.True(t, FlowMatch(ctx, flow1, flow2))

	// different flags
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1, "cookie": 38268468, "flags": 14},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.True(t, FlowMatch(ctx, flow1, flow2))

	// different match InPort
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1, "cookie": 38268468, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 4),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowMatch(ctx, flow1, flow2))

	// different match Ipv4Dst
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1, "cookie": 38268468, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowMatch(ctx, flow1, flow2))

	// different actions
	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1, "cookie": 38268468, "flags": 12},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			PopVlan(ctx),
			Output(ctx, 1),
		},
	}
	flow2, err = MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.True(t, FlowMatch(ctx, flow1, flow2))
}

func TestFlowMatchesMod(t *testing.T) {
	ctx := context.Background()
	assert.False(t, FlowMatchesMod(ctx, nil, nil))
	fa := &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			Output(ctx, 1),
			Group(ctx, 10),
		},
	}
	flow, err := MkFlowStat(ctx, fa)
	assert.Nil(t, err)
	assert.False(t, FlowMatchesMod(ctx, flow, nil))

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"priority": 500, "table_id": 1},
		MatchFields: []*ofp.OfpOxmOfbField{
			InPort(ctx, 2),
			VlanVid(ctx, uint32(ofp.OfpVlanId_OFPVID_PRESENT)|0),
			VlanPcp(ctx, 0),
			EthType(ctx, 0x800),
			Ipv4Dst(ctx, 0xe00a0a0a),
		},
		Actions: []*ofp.OfpAction{
			PopVlan(ctx),
			Output(ctx, 1),
		},
	}
	flowMod := MkSimpleFlowMod(ctx, ToOfpOxmField(ctx, fa.MatchFields), fa.Actions, fa.Command, fa.KV)
	assert.False(t, FlowMatchesMod(ctx, nil, flowMod))
	assert.False(t, FlowMatchesMod(ctx, flow, flowMod))
	entry, err := FlowStatsEntryFromFlowModMessage(ctx, flowMod)
	assert.Nil(t, err)
	assert.True(t, FlowMatch(ctx, flow, entry))

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"table_id": uint64(ofp.OfpTable_OFPTT_ALL),
			"cookie_mask": 0,
			"out_port":    uint64(ofp.OfpPortNo_OFPP_ANY),
			"out_group":   uint64(ofp.OfpGroup_OFPG_ANY),
		},
	}
	flowMod = MkSimpleFlowMod(ctx, ToOfpOxmField(ctx, fa.MatchFields), fa.Actions, fa.Command, fa.KV)
	assert.True(t, FlowMatchesMod(ctx, flow, flowMod))

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"table_id": 1,
			"cookie_mask": 0,
			"out_port":    uint64(ofp.OfpPortNo_OFPP_ANY),
			"out_group":   uint64(ofp.OfpGroup_OFPG_ANY),
		},
	}
	flowMod = MkSimpleFlowMod(ctx, ToOfpOxmField(ctx, fa.MatchFields), fa.Actions, fa.Command, fa.KV)
	assert.True(t, FlowMatchesMod(ctx, flow, flowMod))

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"table_id": 1,
			"cookie_mask": 0,
			"out_port":    1,
			"out_group":   uint64(ofp.OfpGroup_OFPG_ANY),
		},
	}
	flowMod = MkSimpleFlowMod(ctx, ToOfpOxmField(ctx, fa.MatchFields), fa.Actions, fa.Command, fa.KV)
	assert.True(t, FlowMatchesMod(ctx, flow, flowMod))

	fa = &FlowArgs{
		KV: OfpFlowModArgs{"table_id": 1,
			"cookie_mask": 0,
			"out_port":    1,
			"out_group":   10,
		},
	}
	flowMod = MkSimpleFlowMod(ctx, ToOfpOxmField(ctx, fa.MatchFields), fa.Actions, fa.Command, fa.KV)
	assert.True(t, FlowMatchesMod(ctx, flow, flowMod))
}

func TestIsMulticastIpAddress(t *testing.T) {
	isMcastIp := IsMulticastIp(context.Background(), 3776315393) //225.22.0.1
	assert.True(t, isMcastIp)
	isMcastIp = IsMulticastIp(context.Background(), 3232243777) //192.168.32.65
	assert.True(t, !isMcastIp)
}

func TestConvertToMulticastMac(t *testing.T) {
	mcastIp := uint32(4001431809)                   //238.129.1.1
	expectedMacInBytes := []byte{1, 0, 94, 1, 1, 1} //01:00:5e:01:01:01
	macInBytes := ConvertToMulticastMacBytes(context.Background(), mcastIp)
	assert.True(t, bytes.Equal(macInBytes, expectedMacInBytes))
}
