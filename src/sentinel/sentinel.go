package sentinel

import (
	"errors"
	"fmt"
	"log"
	"net"
	"proxy/src/variables"
	"strings"
	"time"
)

type Sentinel struct {
	REDIS_IP   string
	REDIS_NAME string
	sentinel   *net.TCPAddr
}

func (s *Sentinel) read(conn *net.TCPConn) (string, error) {
	err := s.write(conn)

	buffer := make([]byte, 256)

	_, err = conn.Read(buffer)

	parts := strings.Split(string(buffer), "\r\n")

	if err != nil || len(parts) < 5 {
		return "", errors.New("failed to get sentinel")
	}
	return fmt.Sprintf("%s:%s", parts[2], parts[4]), err
}

func (s *Sentinel) write(conn *net.TCPConn) error {
	_, err := conn.Write([]byte(fmt.Sprintf(variables.SENTINEL_COMMAND, s.REDIS_NAME)))
	return err
}

func (s *Sentinel) checkMaster(ip string) {
	if len(s.REDIS_IP) == 0 && s.REDIS_IP != ip {
		log.Println("new redis master", ip)
		s.REDIS_IP = ip
	}
}

func (s *Sentinel) getMaster(conn *net.TCPConn) {
	for {
		if ip, err := s.read(conn); err == nil {
			s.checkMaster(ip)
		} else {
			go s.connect()
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func (s *Sentinel) connect() {
	if conn, err := net.DialTCP("tcp", nil, s.sentinel); err == nil {
		s.getMaster(conn)
	} else {
		time.Sleep(1 * time.Second)
		log.Println("wailed to resolve sentinel", s.sentinel.String())
		s.init(s.sentinel.String())
		go s.connect()
	}
}

func (s *Sentinel) init(sentinel string) {
	adr, err := net.ResolveTCPAddr("tcp", sentinel)
	if err != nil {
		log.Fatal("Failed to resolve sentinel address", err)
	}

	s.sentinel = adr
}

func (s *Sentinel) Start(sentinel string) {
	log.Println("starting sentinel")
	s.init(sentinel)
	s.connect()
}

func GetSentinel(name string) *Sentinel {
	return &Sentinel{REDIS_NAME: name}
}
