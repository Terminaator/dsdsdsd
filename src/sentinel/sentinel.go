package sentinel

import (
	"fmt"
	"io"
	"log"
	"net"
	"proxy/src/resp"
	"proxy/src/socket/client"
	"proxy/src/variables"
	"strings"
	"time"
)

type SentinelInterface interface {
	Connect()
	GetRedis() string
}

type Sentinel struct {
	client     client.ClientInterface
	REDIS_ADR  string
	REDIS_NAME string
}

func (s *Sentinel) reader(reader io.Reader) error {
	buf, err := resp.NewReader(reader).ReadObject()

	if err != nil {
		return err
	}

	parts := strings.Split(string(buf), "\r\n")

	if len(parts) > 4 {
		s.newMaster(fmt.Sprintf("%s:%s", parts[2], parts[4]))
	}

	return err
}

func (s *Sentinel) newMaster(adr string) {
	if len(s.REDIS_ADR) == 0 || adr != s.REDIS_ADR {
		s.REDIS_ADR = adr
		log.Println("New sentinel master", s.REDIS_ADR, s.REDIS_NAME, s.client.GetAdr())
	}
}

func (s *Sentinel) write(writer io.Writer) error {
	_, err := writer.Write([]byte(fmt.Sprintf(variables.SENTINEL_COMMAND, s.REDIS_NAME)))
	return err
}

func (s *Sentinel) getMaster(conn *net.TCPConn) {
	for {
		time.Sleep(1 * time.Second)

		if err := s.write(conn); err == nil {
			if err := s.reader(conn); err != nil {
				go s.Connect()
				break
			}
		} else {
			go s.Connect()
			break
		}
	}
}

func (s *Sentinel) Connect() {
	log.Println("Starting sentinel")
	s.getMaster(s.client.Dial())
}

func (s *Sentinel) GetRedis() string {
	return s.REDIS_ADR
}

func GetSentinel(name string, adr string) SentinelInterface {
	return &Sentinel{REDIS_NAME: name, client: client.GetClient(adr)}
}
