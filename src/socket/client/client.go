package client

import (
	"log"
	"net"
	"time"
)

type ClientInterface interface {
	Dial() *net.TCPConn
	GetAdr() string
}

type Client struct {
	Adr string
}

func (c *Client) resolve() *net.TCPAddr {
	adr, err := net.ResolveTCPAddr("tcp", c.Adr)

	if err != nil {
		log.Fatal("Failed to resolve sentinel address", err)
	}

	return adr
}

func (c *Client) Dial() *net.TCPConn {
	if conn, err := net.DialTCP("tcp", nil, c.resolve()); err == nil {
		return conn
	}

	time.Sleep(1 * time.Second)
	return c.Dial()
}

func (c *Client) GetAdr() string {
	return c.Adr
}

func GetClient(adr string) ClientInterface {
	return &Client{Adr: adr}
}
