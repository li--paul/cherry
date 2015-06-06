/*
 * Cherry - An OpenFlow Controller
 *
 * Copyright (C) 2015 Samjung Data Service Co., Ltd.,
 * Kitae Kim <superkkt@sds.co.kr>
 */

package of10

import (
	"encoding/binary"
	"git.sds.co.kr/cherry.git/cherryd/openflow"
)

type FeaturesRequest struct {
	openflow.Message
}

func NewFeaturesRequest(xid uint32) openflow.FeaturesRequest {
	return &FeaturesRequest{
		Message: openflow.NewMessage(openflow.OF10_VERSION, OFPT_FEATURES_REQUEST, xid),
	}
}

func (r *FeaturesRequest) MarshalBinary() ([]byte, error) {
	return r.Message.MarshalBinary()
}

type FeaturesReply struct {
	openflow.Message
	dpid         uint64
	numBuffers   uint32
	numTables    uint8
	capabilities uint32
	actions      uint32
	ports        []openflow.Port
}

func (r FeaturesReply) DPID() uint64 {
	return r.dpid
}

func (r FeaturesReply) NumBuffers() uint32 {
	return r.numBuffers
}

func (r FeaturesReply) NumTables() uint8 {
	return r.numTables
}

func (r FeaturesReply) Capabilities() uint32 {
	return r.capabilities
}

func (r FeaturesReply) Actions() uint32 {
	return r.actions
}

func (r FeaturesReply) Ports() []openflow.Port {
	return r.ports
}

func (r FeaturesReply) AuxID() uint8 {
	// OpenFlow 1.0 does not have auxilary ID
	return 0
}

func (r *FeaturesReply) UnmarshalBinary(data []byte) error {
	if err := r.Message.UnmarshalBinary(data); err != nil {
		return err
	}

	payload := r.Payload()
	if payload == nil || len(payload) < 24 {
		return openflow.ErrInvalidPacketLength
	}
	r.dpid = binary.BigEndian.Uint64(payload[0:8])
	r.numBuffers = binary.BigEndian.Uint32(payload[8:12])
	r.numTables = payload[12]
	r.capabilities = binary.BigEndian.Uint32(payload[16:20])
	r.actions = binary.BigEndian.Uint32(payload[20:24])

	nPorts := (len(payload) - 24) / 48
	if nPorts == 0 {
		return nil
	}
	r.ports = make([]openflow.Port, nPorts)
	for i := 0; i < nPorts; i++ {
		buf := payload[24+i*48:]
		r.ports[i] = new(Port)
		if err := r.ports[i].UnmarshalBinary(buf[0:48]); err != nil {
			return err
		}
	}

	return nil
}
