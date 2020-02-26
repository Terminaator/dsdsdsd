package main

import (
	"log"
	"proxy/src/proxy"
	"proxy/src/sentinel"
)

func main() {
	log.Println("starting proxy")

	sentinel := sentinel.GetSentinel("mymaster", ":26379")

	go sentinel.Connect()

	proxy := proxy.GetProxy(":9999", sentinel)

	proxy.Start()

}
