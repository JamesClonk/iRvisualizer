package web

import (
	"crypto/subtle"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	visualizerErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "irvisualizer_errors_total",
		Help: "Total errors from iRvisualizer, should be a rate of 0.",
	})
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
	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	// fake index html
	r.HandleFunc("/", h.index)
	r.HandleFunc("/season/", h.index)
	r.HandleFunc("/season/{seasonID}/", h.indexHeatmap)
	r.HandleFunc("/season/{seasonID}/week/", h.index)
	r.HandleFunc("/season/{seasonID}/week/{week}/", h.indexHeatmap)
	r.HandleFunc("/season/{seasonID}/week/{week}/top/", h.indexTop)
	// fake banner
	r.HandleFunc("/banner.png", h.banner)

	// data export
	r.HandleFunc("/series", h.series)
	r.HandleFunc("/series_json", h.seriesJson)
	r.HandleFunc("/series/{seriesID}", h.seriesWeeklyExport) // backwards-compatible endpoint
	r.HandleFunc("/series/{seriesID}/weekly", h.seriesWeeklyExport)
	r.HandleFunc("/series/{seriesID}/week", h.seriesWeeklyExport)
	r.HandleFunc("/series/{seriesID}/season", h.seriesSeasonExport)
	r.HandleFunc("/series/{seriesID}/seasonal", h.seriesSeasonExport)

	// dynamic ranking/standings
	r.HandleFunc("/season/{seasonID}/standings.png", h.ranking)
	r.HandleFunc("/season/{seasonID}/standing.png", h.ranking)
	r.HandleFunc("/season/{seasonID}/rankings.png", h.ranking)
	r.HandleFunc("/season/{seasonID}/ranking.png", h.ranking)
	r.HandleFunc("/season/{seasonID}/oval_standings.png", h.ovalRanking)
	r.HandleFunc("/season/{seasonID}/oval_standing.png", h.ovalRanking)
	r.HandleFunc("/season/{seasonID}/oval_rankings.png", h.ovalRanking)
	r.HandleFunc("/season/{seasonID}/oval_ranking.png", h.ovalRanking)

	// dynamic heatmap
	r.HandleFunc("/season/{seasonID}/week/{week}/heatmap.png", h.weeklyHeatmap)
	r.HandleFunc("/season/{seasonID}/heatmap.png", h.seasonalHeatmap)

	// dynamic scores
	r.HandleFunc("/season/{seasonID}/week/{week}/top/scores.png", h.weeklyTopScores)
	r.HandleFunc("/season/{seasonID}/week/{week}/top/racers.png", h.weeklyTopRacers)
	r.HandleFunc("/season/{seasonID}/week/{week}/top/laps.png", h.weeklyTopLaps)
	r.HandleFunc("/season/{seasonID}/week/{week}/top/safety.png", h.weeklyTopSafety)

	// dynamic driver summaries
	r.HandleFunc("/season/{seasonID}/summary.png", h.seasonSummary)
	r.HandleFunc("/season/{seasonID}/week/{week}/summary.png", h.weeklySummary)

	// dynamic laptime chart
	r.HandleFunc("/season/{seasonID}/week/{week}/laptimes.png", h.weeklyLaptimes)

	// catch-all
	r.PathPrefix("/").HandlerFunc(h.index)

	// add logging
	r.Use(logging)

	// add cache-control headers
	r.Use(caching)

	return r
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.RequestURI, ".png") {
			log.Debugf("received request: %v; %v; %v; %v;", req.UserAgent(), req.Proto, req.Method, req.RequestURI)
		}
		next.ServeHTTP(rw, req)
	})
}

func caching(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Cache-Control", "private, max-age=900, s-maxage=900")
		next.ServeHTTP(rw, req)
	})
}

func (h *Handler) failure(rw http.ResponseWriter, req *http.Request, err error) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(500)
	_, _ = rw.Write([]byte(fmt.Sprintf(`{ "error": "%v" }`, err.Error())))
	visualizerErrors.Inc()
}

func (h *Handler) health(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	_, _ = rw.Write([]byte(`{ "status": "ok" }`))
}

func (h *Handler) verifyBasicAuth(rw http.ResponseWriter, req *http.Request) bool {
	user, pw, ok := req.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(h.Username)) != 1 || subtle.ConstantTimeCompare([]byte(pw), []byte(h.Password)) != 1 {
		rw.Header().Set("WWW-Authenticate", `Basic realm="iRvisualizer"`)
		rw.WriteHeader(401)
		_, _ = rw.Write([]byte("Unauthorized"))
		return false
	}
	return true
}

func (h *Handler) banner(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "image/png")
	rw.WriteHeader(200)
	http.ServeFile(rw, req, "public/banner.png")
}

func (h *Handler) index(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(200)

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
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(200)

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
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(200)

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
