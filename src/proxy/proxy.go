package proxy

import (
	"log"
	"net"
	"proxy/src/clients"
	"proxy/src/proxy/connection"
	"proxy/src/sentinel"
)

type ProxyInterface interface {
	Start(string)
}

type Proxy struct {
	listener *net.TCPListener
	sentinel *sentinel.Sentinel
	clients  *clients.Clients
}

func (p *Proxy) init(listen string) *net.TCPListener {
	laddr, err := net.ResolveTCPAddr("tcp", listen)
	if err != nil {
		log.Fatal("Failed to resolve local address", err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	return listener
}

func (p *Proxy) connection(conn *net.TCPConn) {
	connection := connection.GetConnection(
		p.sentinel,
		conn,
		conn.RemoteAddr().(*net.TCPAddr).String(),
		p.clients)

	go connection.InPipe()
}

func (p *Proxy) listen() {
	log.Println("listening proxy", p.listener.Addr().String())
	for {
		if conn, err := p.listener.AcceptTCP(); err == nil {
			go p.connection(conn)
		}
	}
}

func (p *Proxy) Start(listen string) {
	p.listener = p.init(listen)

	p.listen()
}

func GetProxy(sentinel *sentinel.Sentinel, clients *clients.Clients) ProxyInterface {
	return &Proxy{sentinel: sentinel, clients: clients}
}
