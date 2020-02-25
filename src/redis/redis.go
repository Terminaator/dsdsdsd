package redis

import (
	"errors"
	"log"
	"net"
	"proxy/src/sentinel"
	"proxy/src/util"
	"proxy/src/variables"
	"time"
)

type Redis struct {
	conn     *net.TCPConn
	Sentinel *sentinel.Sentinel
	Host_ip  string
}

func (r *Redis) Close() {
	if r.conn != nil {
		defer r.conn.Close()
		r.conn.Write(variables.REDIS_QUIT)
	}
}

func (r *Redis) read() []byte {
	if buf, err := util.NewReader(r.conn).ReadObject(); err == nil {
		return buf
	} else {
		return variables.ERROR_READ_REDIS
	}
}

func (r *Redis) cmd(out []byte) error {
	if r.conn != nil {
		_, err := r.conn.Write(out)
		return err
	} else {
		return errors.New("conn is null")
	}
}

func (r *Redis) Do(out []byte) []byte {
	return r.do(0, out)
}

func (r *Redis) do(t int, out []byte) []byte {
	if t == 10 {
		return variables.ERROR_TIMEOUT_REDIS
	} else if r.conn == nil {
		if t != 0 {
			time.Sleep(1 * time.Second)
		}
		r.conn = r.connect()
	} else {
		if err := r.cmd(out); err == nil {
			return r.read()
		} else {
			r.conn = nil
		}
	}
	return r.do(t+1, out)
}

func (r *Redis) connect() *net.TCPConn {
	log.Println("making redis connection", r.Host_ip, "sentinel", r.Sentinel.REDIS_IP)
	var conn *net.TCPConn
	addr, _ := net.ResolveTCPAddr("tcp", r.Sentinel.REDIS_IP)

	if c, err := net.DialTCP("tcp", nil, addr); err == nil {
		conn = c
	}

	return conn
}
