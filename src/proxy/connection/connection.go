package connection

import (
	"io"
	"log"
	"proxy/src/clients"
	"proxy/src/redis"
	"proxy/src/sentinel"
	"proxy/src/util"
)

type ConnectionInterface interface {
	InPipe()
}

type Connection struct {
	host_ip string
	redis   redis.Redis
	rwc     io.ReadWriteCloser
	clients *clients.Clients
}

func (c *Connection) close() {
	log.Println("connection closed", c.host_ip)
	c.redis.Close()
	c.rwc.Close()
}

func (c *Connection) do(in []byte) []byte {
	return c.redis.Do(in)
}

func (c *Connection) checkOut(out string) {

}

func (c *Connection) outPipe(out []byte) {
	log.Println("message", out, len(out), "to", c.host_ip)
	go c.checkOut(string(out))
	c.rwc.Write(out)
	c.InPipe()
}

func (c *Connection) InPipe() {
	if in, err := util.NewReader(c.rwc).ReadObject(); err == nil {
		log.Println("message", in, "from", c.host_ip)
		c.outPipe(c.do(in))
	} else {
		c.close()
	}
}

func GetConnection(sentinel *sentinel.Sentinel, rwc io.ReadWriteCloser, ip string, clients *clients.Clients) ConnectionInterface {
	log.Println("new connection from", ip)
	return &Connection{
		host_ip: ip,
		rwc:     rwc,
		redis:   redis.Redis{Host_ip: ip, Sentinel: sentinel},
		clients: clients}
}
