package ohrad

import (
	"log"
	"net"
	"sync"
	"time"
)

const (
	VersionMask = 7 << 3
	ModeMask    = 7 << 0
	ModeServer  = 4
	ModeClient  = 3
	ModeSynAct  = 1
	ModeSynPas  = 2
	Jan1970     = 2208988800
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
			log.Println(err)
			continue
		}

		srv.inFlight.Add(1)
		go srv.serve(conn, addr, query)

	}
}

func (srv *Server) serve(conn *net.UDPConn, clientAddr *net.UDPAddr, query NtpMsg) {
	defer srv.inFlight.Done()

	var reply NtpMsg

	log.Println(query)

	// todo: move right after readfromudp
	rectime := NtpLong{
		IntParl: unix2ntp(time.Now().Unix()),
		//Fractionl: unix2ntp(time.Now().Nanosecond()),
		Fractionl: 0,
	}

	// Header
	//reply.Status = 3 << 6
	reply.Status = 0
	reply.Status |= (query.Status & VersionMask)
	if (query.Status & ModeMask) == ModeClient {
		reply.Status |= ModeServer
	} else if (query.Status & ModeMask) == ModeSynAct {
		reply.Status |= ModeSynPas
	} else {
		log.Println("Invalid msg")
		return
	}

	// Body
	reply.Stratum = 3
	reply.Ppoll = query.Ppoll
	reply.Precision = 0
	reply.Rectime = rectime
	reply.Reftime = rectime
	reply.Orgtime = query.Xmttime
	reply.Rootdelay = NtpShort{0, 0}
	reply.Refid = 0xc8a007ba
	reply.Xmttime = NtpLong{
		IntParl: unix2ntp(time.Now().Unix()),
		//Fractionl: unix2ntp(time.Now().Nanosecond()),
		Fractionl: 0,
	}

	log.Println(reply)

	SendNtpMsg(conn, clientAddr, &reply)

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
