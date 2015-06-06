/*
 * Cherry - An OpenFlow Controller
 *
 * Copyright (C) 2015 Samjung Data Service Co., Ltd.,
 * Kitae Kim <superkkt@sds.co.kr>
 */

package of13

import (
	"encoding/binary"
	"git.sds.co.kr/cherry.git/cherryd/openflow"
	"strings"
)

type DescRequest struct {
	openflow.Message
}

func NewDescRequest(xid uint32) openflow.DescRequest {
	return &DescRequest{
		Message: openflow.NewMessage(openflow.OF13_VERSION, OFPT_MULTIPART_REQUEST, xid),
	}
}

func (r *DescRequest) MarshalBinary() ([]byte, error) {
	v := make([]byte, 8)
	// Multipart description request
	binary.BigEndian.PutUint16(v[0:2], OFPMP_DESC)
	// No flags and body
	r.SetPayload(v)

	return r.Message.MarshalBinary()
}

type DescReply struct {
	openflow.Message
	manufacturer string
	hardware     string
	software     string
	serial       string
	description  string
}

func (r DescReply) Manufacturer() string {
	return r.manufacturer
}

func (r DescReply) Hardware() string {
	return r.hardware
}

func (r DescReply) Software() string {
	return r.software
}

func (r DescReply) Serial() string {
	return r.serial
}

func (r DescReply) Description() string {
	return r.description
}

func (r *DescReply) UnmarshalBinary(data []byte) error {
	if err := r.Message.UnmarshalBinary(data); err != nil {
		return err
	}

	payload := r.Payload()
	if payload == nil || len(payload) < 1064 {
		return openflow.ErrInvalidPacketLength
	}
	r.manufacturer = strings.TrimRight(string(payload[8:264]), "\x00")
	r.hardware = strings.TrimRight(string(payload[264:520]), "\x00")
	r.software = strings.TrimRight(string(payload[520:776]), "\x00")
	r.serial = strings.TrimRight(string(payload[776:808]), "\x00")
	r.description = strings.TrimRight(string(payload[808:1064]), "\x00")

	return nil
}
