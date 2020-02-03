package web

import (
	"crypto/subtle"
	"fmt"
	"html/template"
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
	return router(h)
}

func router(h *Handler) *mux.Router {
	// mux router
	r := mux.NewRouter()
	r.PathPrefix("/health").HandlerFunc(h.health)

	// fake index html
	r.HandleFunc("/", h.index).Methods("GET")
	r.HandleFunc("/season/", h.index).Methods("GET")
	r.HandleFunc("/season/{seasonID}/", h.indexHeatmap).Methods("GET")
	r.HandleFunc("/season/{seasonID}/week/", h.index).Methods("GET")
	r.HandleFunc("/season/{seasonID}/week/{week}/", h.indexHeatmap).Methods("GET")
	r.HandleFunc("/season/{seasonID}/week/{week}/top/", h.indexTop).Methods("GET")

	// dynamic ranking/standings
	r.HandleFunc("/season/{seasonID}/standings.png", h.ranking).Methods("GET")
	r.HandleFunc("/season/{seasonID}/standing.png", h.ranking).Methods("GET")
	r.HandleFunc("/season/{seasonID}/rankings.png", h.ranking).Methods("GET")
	r.HandleFunc("/season/{seasonID}/ranking.png", h.ranking).Methods("GET")

	// dynamic heatmap
	r.HandleFunc("/season/{seasonID}/week/{week}/heatmap.png", h.weeklyHeatmap).Methods("GET")
	r.HandleFunc("/season/{seasonID}/heatmap.png", h.seasonalHeatmap).Methods("GET")

	// dynamic scores
	r.HandleFunc("/season/{seasonID}/week/{week}/top/scores.png", h.weeklyTopScores).Methods("GET")
	r.HandleFunc("/season/{seasonID}/week/{week}/top/racers.png", h.weeklyTopRacers).Methods("GET")
	r.HandleFunc("/season/{seasonID}/week/{week}/top/laps.png", h.weeklyTopLaps).Methods("GET")
	r.HandleFunc("/season/{seasonID}/week/{week}/top/safety.png", h.weeklyTopSafety).Methods("GET")

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

func (h *Handler) index(rw http.ResponseWriter, req *http.Request) {
	index, _ := template.New("index").Parse(`
<html>
<head>
<title>Statistics</title>
<body>
nothing here...
</body>
</html>
	`)
	if err := index.ExecuteTemplate(rw, "index", nil); err != nil {
		h.failure(rw, req, err)
	}
}

func (h *Handler) indexHeatmap(rw http.ResponseWriter, req *http.Request) {
	index, _ := template.New("index").Parse(`
<html>
<head>
<title>Heatmap</title>
<body>
<img src="heatmap.png"/><br/>
</body>
</html>
	`)
	if err := index.ExecuteTemplate(rw, "index", nil); err != nil {
		h.failure(rw, req, err)
	}
}

func (h *Handler) indexTop(rw http.ResponseWriter, req *http.Request) {
	index, _ := template.New("index").Parse(`
<html>
<head>
<title>Top Racers</title>
<body>
<img src="scores.png"/><br/>
<img src="racers.png"/><br/>
<img src="safety.png"/><br/>
<img src="laps.png"/><br/>
</body>
</html>
	`)
	if err := index.ExecuteTemplate(rw, "index", nil); err != nil {
		h.failure(rw, req, err)
	}
}
