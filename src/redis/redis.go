package redis

import (
	"net"
	"proxy/src/resp"
	"proxy/src/sentinel"
	"proxy/src/variables"
)

type RedisInterface interface {
	Do([]byte) []byte
	Close()
}

type Redis struct {
	conn     *net.TCPConn
	sentinel sentinel.SentinelInterface
}

func (r *Redis) read() []byte {
	if buf, err := resp.NewReader(r.conn).ReadObject(); err == nil {
		return buf
	}
	return variables.ERROR_READ_REDIS
}

func (r *Redis) cmd(in []byte) error {
	if r.conn != nil {
		_, err := r.conn.Write(in)
		return err
	}

	return variables.ErrWriting
}

func (r *Redis) Do(in []byte) []byte {
	return r.do(in, 0)
}

func (r *Redis) do(in []byte, timeout int) []byte {
	if r.conn == nil {
		r.connect()
	}

	if err := r.cmd(in); err == nil {
		return r.read()
	}

	r.conn = nil

	if timeout == 10 {
		return variables.ERROR_TIMEOUT_REDIS
	}

	return r.do(in, timeout+1)

}

func (r *Redis) Close() {
	if r.conn != nil {
		r.conn.Write(variables.REDIS_QUIT)
		r.conn.Close()
	}
}

func (r *Redis) connect() {
	addr, _ := net.ResolveTCPAddr("tcp", r.sentinel.GetRedis())

	if c, err := net.DialTCP("tcp", nil, addr); err == nil {
		r.conn = c
	}
}

func GetRedis(sentinel sentinel.SentinelInterface) RedisInterface {
	return &Redis{sentinel: sentinel}
}
