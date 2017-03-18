package ohrad

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"syscall"
)

const (
	NtpMsgSize     int = 48 // no auth
	Jan1970            = 2208988800
	NanosPerSecond     = 1000 * 1000 * 1000
)

type NtpLong struct {
	IntParl   uint32
	Fractionl uint32
}

type NtpShort struct {
	IntParts  uint16
	Fractions uint16
}

type NtpMsg struct {
	Status     uint8
	Stratum    uint8
	Ppoll      uint8
	Precision  int8
	Rootdelay  NtpShort
	Dispersion NtpShort
	Refid      uint32
	Reftime    NtpLong
	Orgtime    NtpLong
	Rectime    NtpLong
	Xmttime    NtpLong
}

func (m *NtpMsg) Bytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, m)
	return buf.Bytes()
}

func NewMsg(buf []byte) *NtpMsg {
	var query NtpMsg
	binary.Read(bytes.NewReader(buf), binary.BigEndian, &query)
	return &query
}

func GetNtpMsg(conn *net.UDPConn) (*NtpMsg, *net.UDPAddr, error) {
	// FIXME: only handling no-auth messages for now
	buf := make([]byte, NtpMsgSize)
	n, addr, err := conn.ReadFromUDP(buf)

	rectime := getTimeNow()

	if n != NtpMsgSize {
		Log.Debug("Invalid msg size)")
	}

	query := NewMsg(buf[0:n])
	query.Rectime = rectime

	return query, addr, err
}

func SendNtpMsg(conn *net.UDPConn, clientAddr *net.UDPAddr, msg *NtpMsg) {
	msgBytes := msg.Bytes()
	n, nerr := conn.WriteToUDP(msgBytes, clientAddr)

	if errno, ok := nerr.(syscall.Errno); ok {
		if (errno == syscall.ENOBUFS) || (errno == syscall.EHOSTUNREACH) || (errno == syscall.ENETDOWN) || (errno == syscall.EHOSTDOWN) {
			return
		}
		Log.Debug(fmt.Sprintf("WriteToUDP: %s", errno))
		return
	}

	if n != len(msgBytes) {
		Log.Notice("Sent msg has a different size")
	}

	return
}
