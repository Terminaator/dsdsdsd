package main

import (
	"log"
	"os"
	"proxy/src/api/router"
	"proxy/src/clients"
	"proxy/src/proxy"
	"proxy/src/sentinel"
)

func main() {
	log.Println("starting proxy")

	CLIENTS_FILE_PATH, SENTINEL_MASTER_NAME, SENTINEL_ADR, API_TOKEN, PROXY_PORT, PROXY_API_PORT := getEnvValues()

	sentinel := sentinel.GetSentinel(SENTINEL_MASTER_NAME)
	go sentinel.Start(SENTINEL_ADR)

	clients := clients.Clients{Sentinel: sentinel, State: clients.Start}
	go clients.Start(CLIENTS_FILE_PATH)

	router := router.Router{}
	go router.Start(PROXY_API_PORT, sentinel, API_TOKEN, &clients)

	proxy := proxy.GetProxy(sentinel, &clients)

	proxy.Start(PROXY_PORT)
}

func getEnvValues() (string, string, string, string, string, string) {
	return getEnv("CLIENTS_FILE_PATH", "./conf/init.json"),
		getEnv("SENTINEL_MASTER_NAME", "mymaster"),
		getEnv("SENTINEL_IP", "127.0.0.1") + ":" + getEnv("SENTINEL_PORT", "26379"),
		getEnv("API_TOKEN", "d29bbdfd9f7c2b46d142590330f28ef9029da92a83c947b57924504fd7f4abc092a550eb868f3ebf2d7f152690de26e975cf991e7b5d47bdeabf8990c89d09ed32ee8e18ca7ae62d13fd302cfc2683c5e39e398c38cf2b0e82f7ff764b30a8af587b651a"),
		getEnv("PROXY_PORT", ":9999"),
		getEnv("PROXY_API_PORT", ":8080")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
