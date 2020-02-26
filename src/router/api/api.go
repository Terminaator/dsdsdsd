package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"proxy/src/redis"
	"proxy/src/sentinel"
	"strings"

	"github.com/gorilla/mux"
)

type Api struct {
	Router   *mux.Router
	sentinel sentinel.SentinelInterface
}

type Respond struct {
	Respond string
}

func (a *Api) Start() {
	a.Router.HandleFunc("/building", a.getBuilding)
	a.Router.HandleFunc("/utilitybuilding", a.getUtilityBuilding)
	a.Router.HandleFunc("/procedure", a.getProcedure)
	a.Router.HandleFunc("/document/{doty}", a.getDocument)
	a.Router.HandleFunc("/readiness", a.getReadiness)
}

func (a *Api) getBuilding(w http.ResponseWriter, _ *http.Request) {
	a.response(&w, fmt.Sprintf(variables.REDIS_ONE_COMMAND, len(variables.BUILDING_CODE), variables.BUILDING_CODE))
}

func (a *Api) getUtilityBuilding(w http.ResponseWriter, _ *http.Request) {
	a.response(&w, fmt.Sprintf(variables.REDIS_ONE_COMMAND, len(variables.UTILITY_BUILDING_CODE), variables.UTILITY_BUILDING_CODE))
}

func (a *Api) getProcedure(w http.ResponseWriter, _ *http.Request) {
	a.response(&w, fmt.Sprintf(variables.REDIS_ONE_COMMAND, len(variables.PROCEDURE_CODE), variables.PROCEDURE_CODE))
}

func (a *Api) getDocument(w http.ResponseWriter, r *http.Request) {
	a.response(&w, fmt.Sprintf(variables.REDIS_TWO_COMMAND, len(variables.DOCUMENT_CODE), variables.DOCUMENT_CODE, len(mux.Vars(r)["doty"]), mux.Vars(r)["doty"]))
}

func (a *Api) getReadiness(w http.ResponseWriter, _ *http.Request) {
	if len(a.Sentinel.REDIS_IP) != 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, "")
}

func (a *Api) checkOut(out string) {
}

func (a *Api) doRedis(command []byte) (string, error) {
	redis := redis.Redis{Sentinel: a.Sentinel, Host_ip: "api"}

	r := redis.Do(command)

	out := string(r)
	defer redis.Close()

	if out[0] == '-' {
		go a.checkOut(out)
		return strings.Trim(out[1:], string(variables.REDIS_END)), errors.New("err")
	} else {
		return strings.Trim(out[1:], string(variables.REDIS_END)), nil
	}

}

func (a *Api) response(w *http.ResponseWriter, command string) {
	log.Println("api redis", []byte(command))

	r, err := a.doRedis([]byte(command))

	res := Respond{r}

	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(*w).Encode(res)

}
