package clients

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"proxy/src/redis"
	"proxy/src/sentinel"
	"reflect"
	"strconv"
	"time"
)

type ClientsState int

const (
	Start ClientsState = iota
	Old
	New
	Init
)

type Clients struct {
	list     ClientList
	redis_ip string
	Sentinel *sentinel.Sentinel
	State    ClientsState
}

type ClientList struct {
	File    string
	Clients []string
}

func (i *ClientList) readFile() {
	jsonFile, err := os.Open(i.File)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, i)
}

func (c *Clients) getClients(path string) {
	c.list = ClientList{File: path}
	c.list.readFile()

	log.Println(c.list)
}

func (c *Clients) makeClientRequests(clientUrl string) *map[string]interface{} {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(clientUrl)

	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Error occured when getting response from client or statuscode is not 200")
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	return &result

}

func (c *Clients) addBiggestValuesIntoMap(final *map[string]interface{}, client *map[string]interface{}) {
	for key, value := range *client {
		switch v := value.(type) {
		case map[string]interface{}:
			m := c.convertMap(&v)
			if finalValue, ok := (*final)[key]; ok {
				finalSubValue := finalValue.(map[string]interface{})
				c.addBiggestValuesIntoMap(&finalSubValue, m)
			} else {
				(*final)[key] = *m
			}
		case float64:
			c.addBiggestValueIntoMap(final, &key, int(v))
		case int:
			c.addBiggestValueIntoMap(final, &key, v)
		default:
			log.Fatal("Wrong type")
		}
	}
}

func (c *Clients) convertMap(m *map[string]interface{}) *map[string]interface{} {
	var r = make(map[string]interface{})
	for k, v := range *m {
		r[k] = int(v.(float64))
	}
	return &r
}

func (c *Clients) addBiggestValueIntoMap(final *map[string]interface{}, key *string, value int) {
	if finalValue, ok := (*final)[*key]; ok {
		if finalValue.(int) < value {
			(*final)[*key] = value
		}
	} else {
		(*final)[*key] = value
	}
}

func (c *Clients) getValuesFromClients() *map[string]interface{} {
	final := make(map[string]interface{})
	for _, client := range c.list.Clients {
		log.Println("getting values from", client)
		c.addBiggestValuesIntoMap(&final, c.makeClientRequests(client))
	}
	log.Println(final)
	return &final
}

func (c *Clients) doRedis(command []byte) {
	log.Println("clients add value", command)
	redis := redis.Redis{Sentinel: c.Sentinel, Host_ip: "clients"}

	out := redis.Do(command)

	log.Println("client out", out)
	defer redis.Close()
}

func (c *Clients) addValuesIntoRedis(m *map[string]interface{}) {
	for k, v := range *m {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			for k2, v2 := range v.(map[string]interface{}) {
				command := fmt.Sprintf("*4\r\n$4\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "HSET", len(k), k, len(k2), k2, len(strconv.Itoa(v2.(int))), v2.(int))
				c.doRedis([]byte(command))
			}
		} else {
			command := fmt.Sprintf("*3\r\n$3\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "SET", len(k), k, len(strconv.Itoa(v.(int))), v.(int))
			c.doRedis([]byte(command))
		}
	}
}

func (c *Clients) checkMaster(ip string) {
	if len(ip) != 0 && c.redis_ip != ip {
		c.redis_ip = ip
		c.State = New
	}
}

func (c *Clients) doInit() {
	if c.State == New {
		c.State = Init
		c.addValuesIntoRedis(c.getValuesFromClients())
		c.State = Old
	}
}

func (c *Clients) Init() {
	if c.State == Old {
		c.State = New
	}
}

func (c *Clients) Start(path string) {
	log.Println("starting clients")

	c.getClients(path)

	for {
		//c.checkMaster(c.Sentinel.REDIS_IP)
		c.doInit()
		time.Sleep(1 * time.Second)
	}
}
