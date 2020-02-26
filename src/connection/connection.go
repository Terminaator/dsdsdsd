package connection

import (
	"io"
	"log"
	"proxy/src/redis"
	"proxy/src/resp"
	"proxy/src/sentinel"
)

type ConnectionInterface interface {
	InPipe()
}

type Connection struct {
	rwc   *io.ReadWriteCloser
	adr   *string
	redis redis.RedisInterface
}

func (c *Connection) close(err error) {
	log.Println("Connection closed from", *c.adr, "cause", err)
	c.redis.Close()
	(*c.rwc).Close()
}

func (c *Connection) outPipe(out []byte) {
	log.Println("Proxy message out", *c.adr, "message", out)
	(*c.rwc).Write(out)
	c.InPipe()
}

func (c *Connection) redisPipe(in []byte) {
	log.Println("Proxy message from", *c.adr, "message", in)
	c.outPipe(c.redis.Do(in))

}

func (c *Connection) InPipe() {
	if in, err := resp.NewReader(*c.rwc).ReadObject(); err == nil {
		c.redisPipe(in)
	} else {
		c.close(err)
	}
}

func GetConnection(sentinel sentinel.SentinelInterface, rwc io.ReadWriteCloser, adr string) ConnectionInterface {
	return &Connection{rwc: &rwc, adr: &adr, redis: redis.GetRedis(sentinel)}
}
