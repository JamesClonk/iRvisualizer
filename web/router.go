package web

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"sync"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/gorilla/mux"
)

type Handler struct {
	Username string
	Password string
	DB       database.Database
	Mutex    *sync.Mutex
}

func NewRouter(username, password string) *mux.Router {
	// setup database connection
	db := database.NewDatabase(database.NewAdapter())

	// global handler
	h := &Handler{
		Username: username,
		Password: password,
		DB:       db,
		Mutex:    &sync.Mutex{},
	}

	// mux router
	r := mux.NewRouter()
	r.PathPrefix("/health").HandlerFunc(h.health)

	r.HandleFunc("/season/{seasonID}/week/{week}/heatmap.png", h.weeklyHeatmap).Methods("GET")
	r.HandleFunc("/season/{seasonID}/heatmap.png", h.seasonalHeatmap).Methods("GET")

	return r
}

func (h *Handler) failure(rw http.ResponseWriter, req *http.Request, err error) {
	rw.WriteHeader(500)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(fmt.Sprintf(`{ "error": "%v" }`, err.Error())))
}

func (h *Handler) health(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{ "status": "ok" }`))
}

func (h *Handler) verifyBasicAuth(rw http.ResponseWriter, req *http.Request) bool {
	user, pw, ok := req.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(h.Username)) != 1 || subtle.ConstantTimeCompare([]byte(pw), []byte(h.Password)) != 1 {
		rw.Header().Set("WWW-Authenticate", `Basic realm="iRcollector"`)
		rw.WriteHeader(401)
		rw.Write([]byte("Unauthorized"))
		return false
	}
	return true
}
