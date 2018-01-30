package main

import (
	"net"
	"sync"
	"time"

	dns "github.com/miekg/dns"
)

type (
	Conn struct {
		in    chan []byte
		out   chan []byte
		close chan bool
	}

	Listener struct {
		Listener net.Listener
		refcnt   int
	}

	ListenerProtocol int
	BindRequestType  int

	BindRequestHeader struct {
		Protocol ListenerProtocol `msgpack:"prot"`
		BindInfo string           `msgpack:"bind"`
		Timeout  int              `msgpack:"timeout"`
	}

	DNSRequest struct {
		Name string
		IPs  chan []string
	}

	DNSCacheRecord struct {
		ResponseRecords []dns.RR
		LastActivity    time.Time
	}

	ConnectionAcceptHeader struct {
		Extra      bool   `msgpack:"extra"`
		LocalHost  string `msgpack:"lhost"`
		LocalPort  int    `msgpack:"lport"`
		RemoteHost string `msgpack:"rhost"`
		RemotePort int    `msgpack:"rport"`
		Error      string `msgpack:"error"`
	}

	DNSListener struct {
		Conn net.Conn

		Domain string

		DNSCache    map[string]*DNSCacheRecord
		UDPServer   *dns.Server
		TCPServer   *dns.Server
		DNSRequests chan *DNSRequest

		activeLock sync.Mutex
		active     bool
	}

	Daemon struct {
		Addr string

		DNSLock     sync.Mutex
		DNSCheck    sync.Mutex
		DNSListener *DNSListener

		Listeners     map[int]*Listener
		ListenersLock sync.Mutex
	}

	Extra struct {
		Extra bool   `msgpack:"extra"`
		Data  string `msgpack:"data"`
	}

	TLSAcceptorConfig struct {
		CACert string `msgpack:"ca"`
		Cert   string `msgpack:"cert"`
		Key    string `msgpack:"key"`
	}
)

const (
	DNS ListenerProtocol = 0
	TCP ListenerProtocol = iota
	KCP ListenerProtocol = iota
	TLS ListenerProtocol = iota
)
