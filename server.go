package ohrad

import (
	"net"
	"sync"
	"time"
)

const (
	VersionMask = 7 << 3
	ModeMask    = 7 << 0
	ModeServer  = 4
	ModeClient  = 3
	ModeSymAct  = 1 // Symmetric active
	ModeSymPas  = 2 // Symmetric passive
)

type Server struct {
	Addr        string
	ReadTimeout int32
	started     bool
	inFlight    sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		Addr:        ":123",
		ReadTimeout: 5,
		started:     false,
	}
}

func (srv *Server) getReadTimeout() time.Duration {
	return time.Duration(srv.ReadTimeout) * time.Second
}

func (srv *Server) ListenAndServe() error {
	srvAddr, err := net.ResolveUDPAddr("udp", srv.Addr)
	if err != nil {
		return err
	}

	udpConn, err := net.ListenUDP("udp", srvAddr)
	if err != nil {
		return err
	}

	err = srv.serveUDP(udpConn)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) serveUDP(conn *net.UDPConn) error {
	defer conn.Close()

	for {
		query, addr, err := GetNtpMsg(conn)
		if err != nil {
			Log.Debug(err.Error())
			continue
		}

		srv.inFlight.Add(1)
		go srv.serve(conn, addr, query)

	}
}

func (srv *Server) serve(conn *net.UDPConn, clientAddr *net.UDPAddr, query *NtpMsg) {
	defer srv.inFlight.Done()

	var reply NtpMsg

	Log.NtpMsg(&reply)

	// LI (leap indicator) is a 2-bit field. Not supported for [now|ever]
	reply.Status = 0

	// VN is a 3-bit field representing the protocol version
	reply.Status |= (query.Status & VersionMask)

	// Mode is a 3-bit field representin the type of message
	if (query.Status & ModeMask) == ModeClient {
		reply.Status |= ModeServer
	} else if (query.Status & ModeMask) == ModeSymAct {
		reply.Status |= ModeSymPas
	} else {
		Log.Debug("Invalid msg")
		return
	}

	// We aren't fetching the time from other servers, so I'm making
	// the stratum field as higher as possible. 16 is unsynchronized
	reply.Stratum = 15
	reply.Refid = 0
	reply.Precision = 0
	reply.Ppoll = query.Ppoll
	reply.Rectime = query.Rectime
	reply.Reftime = query.Rectime
	reply.Orgtime = query.Xmttime
	reply.Rootdelay = NtpShort{0, 0}
	reply.Xmttime = getTimeNow()

	SendNtpMsg(conn, clientAddr, &reply)

	Log.NtpMsg(&reply)

	return
}

func (srv *Server) Shutdown() error {
	quit := make(chan bool)

	go func() {
		srv.inFlight.Wait()
		quit <- true
	}()

	select {
	case <-time.After(srv.getReadTimeout()):
		return &ErrorTimeout{where: "shutdown"}
	case <-quit:
		return nil

	}

	return nil
}
