package proxy

import (
	"log"
	"net"
	"proxy/src/connection"
	"proxy/src/sentinel"
	"proxy/src/socket/server"
)

type ProxyInterface interface {
	Start()
}

type Proxy struct {
	server   *server.Server
	sentinel sentinel.SentinelInterface
}

func (p *Proxy) connection(conn *net.TCPConn) {
	log.Println("New proxy connection from", p.server.ListenerAdr())
	go connection.GetConnection(p.sentinel, conn, p.server.ListenerAdr()).InPipe()
}

func (p *Proxy) Start() {
	log.Println("Listening proxy", p.server.ListenerAdr())

	for {
		if conn, err := p.server.Listener().AcceptTCP(); err == nil {
			go p.connection(conn)
		} else {
			log.Fatal("Proxy failed", err)
		}
	}
}

func GetProxy(adr string, sentinel sentinel.SentinelInterface) ProxyInterface {
	return &Proxy{sentinel: sentinel, server: &server.Server{Adr: adr}}
}
