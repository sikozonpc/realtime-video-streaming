package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"goproject/auth"
	json2 "goproject/json"
	"goproject/responses"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTP represents auth http service
type HTTP struct {
	svc auth.Service
}

// NewHTTP creates new auth http service
func NewHTTP(svc auth.Service, r *mux.Router) {
	h := HTTP{svc}

	rp := r.PathPrefix("/auth/").Subrouter()
	rp.HandleFunc("/auth/register", json2.SetMiddlewareJSON(h.register)).Methods("POST")
	rp.HandleFunc("/auth/users", json2.SetMiddlewareJSON(h.register)).Methods("GET")
}

func (h *HTTP) register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	creds := new(auth.Credentials)
	err = json.Unmarshal(body, &creds)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = h.svc.Register(*creds)
	if err != nil {
		log.Fatal(err.Error())
	}

	responses.JSON(w, http.StatusOK, "User registered!")
}

func (h *HTTP) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetUsers()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}
