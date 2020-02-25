package router

import (
	"log"
	"net/http"
	"proxy/src/api"
	"proxy/src/api/auth"
	"proxy/src/clients"
	"proxy/src/sentinel"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func (r *Router) Start(port string, sentinel *sentinel.Sentinel, token string, clients *clients.Clients) {
	log.Println("port", port, "starting api")

	r.router = mux.NewRouter().StrictSlash(true)

	auth := auth.Auth{Token: token}

	auth.Middleware(r.router)

	api := api.Api{Sentinel: sentinel, Router: r.router, Clients: clients}

	api.Start()

	go http.ListenAndServe(port, r.router)
}
