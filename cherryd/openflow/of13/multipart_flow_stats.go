/*
 * Cherry - An OpenFlow Controller
 *
 * Copyright (C) 2015 Samjung Data Service Co., Ltd.,
 * Kitae Kim <superkkt@sds.co.kr>
 */

package of13

import (
	"encoding/binary"
	"errors"
	"git.sds.co.kr/cherry.git/cherryd/openflow"
)

type FlowStatsRequest struct {
	openflow.Message
	tableID            uint8
	cookie, cookieMask uint64
	match              openflow.Match
}

func NewFlowStatsRequest(xid uint32) openflow.FlowStatsRequest {
	return &FlowStatsRequest{
		Message: openflow.NewMessage(openflow.OF13_VERSION, OFPT_MULTIPART_REQUEST, xid),
	}
}

func (r FlowStatsRequest) Cookie() uint64 {
	return r.cookie
}

func (r *FlowStatsRequest) SetCookie(cookie uint64) error {
	r.cookie = cookie
	return nil
}

func (r FlowStatsRequest) CookieMask() uint64 {
	return r.cookieMask
}

func (r *FlowStatsRequest) SetCookieMask(mask uint64) error {
	r.cookieMask = mask
	return nil
}

func (r FlowStatsRequest) Match() openflow.Match {
	return r.match
}

func (r *FlowStatsRequest) SetMatch(match openflow.Match) error {
	if match == nil {
		return errors.New("match is nil")
	}
	r.match = match
	return nil
}

func (r FlowStatsRequest) TableID() uint8 {
	return r.tableID
}

// 0xFF means all table
func (r *FlowStatsRequest) SetTableID(id uint8) error {
	r.tableID = id
	return nil
}

func (r *FlowStatsRequest) MarshalBinary() ([]byte, error) {
	v := make([]byte, 40)
	// Flow stats request
	binary.BigEndian.PutUint16(v[0:2], OFPMP_FLOW)
	v[8] = r.tableID
	// v[9:12] is padding
	binary.BigEndian.PutUint32(v[12:16], OFPP_ANY)
	binary.BigEndian.PutUint32(v[16:20], OFPG_ANY)
	// v[20:24] is padding
	binary.BigEndian.PutUint64(v[24:32], r.cookie)
	binary.BigEndian.PutUint64(v[32:40], r.cookieMask)
	match, err := r.match.MarshalBinary()
	if err != nil {
		return nil, err
	}
	v = append(v, match...)
	r.SetPayload(v)

	return r.Message.MarshalBinary()
}

// TODO: Implement FlowStatsReply
