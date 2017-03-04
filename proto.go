package ohrad

import (
	"bytes"
	"encoding/binary"
)

type NtpLong struct {
	IntParl   uint32
	Fractionl uint32
}

type NtpShort struct {
	IntParts  uint16
	Fractions uint16
}

type Msg struct {
	Status     uint8 // status of local clock and leap info
	Stratum    uint8 // Stratum level
	Ppoll      uint8 // poll value
	Precision  int8
	Rootdelay  NtpShort
	Dispersion NtpShort
	Refid      uint32
	Reftime    NtpLong
	Orgtime    NtpLong
	Rectime    NtpLong
	Xmttime    NtpLong
}

func (p *Msg) Bytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, p)
	return buf.Bytes()
}

func NewMsg(buf []byte) Msg {
	var query Msg
	binary.Read(bytes.NewReader(buf), binary.BigEndian, &query)
	return query
}
