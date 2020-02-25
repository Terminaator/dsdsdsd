package redis

import (
	"bytes"
	"log"
	"net"
	"proxy/src/sentinel"
	"testing"
)

var (
	adr string = ":6380"
)

func listener(adr string) *net.TCPListener {
	laddr, _ := net.ResolveTCPAddr("tcp", adr)

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	return listener
}

func mockRedis() func(t *testing.T, listener *net.TCPListener, in []byte, out []byte) {

	return func(t *testing.T, listener *net.TCPListener, in []byte, out []byte) {
		buf := make([]byte, 128)
		if c, err := listener.AcceptTCP(); err == nil {
			c.Read(buf)

			buf = bytes.Trim(buf, "\x00")
			if string(buf) != string(in) {
				t.Error()
			}

			c.Write(out)
		}
	}
}

func makeTestRedis(redis string) Redis {
	sentinel := sentinel.Sentinel{REDIS_IP: redis}
	return Redis{Host_ip: "127.0.0.1:23324", Sentinel: &sentinel}
}

func do(t *testing.T, in []byte, out []byte) []byte {
	mock := mockRedis()
	redis := makeTestRedis(adr)

	listener := listener(adr)
	go mock(t, listener, in, out)

	buf := redis.Do(in)

	listener.Close()

	return buf
}

func TestRedisSet(t *testing.T) {
	t.Log("redis set test")

	in := []byte("*3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\na\r\n")
	out := []byte("+OK\r\n")

	buf := do(t, in, out)

	t.Log(string(buf))

	if string(buf) != string(out) {
		t.Error()
	}

}

func TestRedisHGet(t *testing.T) {
	t.Log("redis hget test")
	in := []byte("*3\r\n$4\r\nhget\r\n$1\r\n2\r\n$1\r\n2\r\n")
	out := []byte("$1\r\n2\r\n")
	buf := do(t, in, out)

	t.Log(string(buf))

	if string(buf) != string(out) {
		t.Error()
	}

}

func TestRedisHSet(t *testing.T) {
	t.Log("redis hset test")
	in := []byte("*4\r\n$4\r\nhset\r\n$1\r\n2\r\n$1\r\n2\r\n$1\r\n2\r\n")
	out := []byte(":1\r\n")

	buf := do(t, in, out)
	t.Log(string(buf))

	if string(buf) != string(out) {
		t.Error()
	}

}

func TestRedisHGetAll(t *testing.T) {
	t.Log("redis hgetall test")
	in := []byte("*2\r\n$7\r\nhgetall\r\n$1\r\n2\r\n")
	out := []byte("*2\r\n$1\r\n2\r\n$1\r\n2\r\n")

	buf := do(t, in, out)
	t.Log(string(buf))

	if string(buf) != string(out) {
		t.Error()
	}

}

func TestRedisTimeout(t *testing.T) {
	t.Log("redis timeout test")
	redis := makeTestRedis("")

	buf := redis.Do([]byte("*3\r\n$4\r\nhget\r\n$1\r\n2\r\n$1\r\n2\r\n"))
	t.Log(string(buf))

	if string(buf) != "-Timeout\r\n" {
		t.Error()
	}

}
