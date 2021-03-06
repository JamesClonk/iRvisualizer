package web

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/stretchr/testify/assert"
)

func Test_Health(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := &Handler{
		Username: "iracing",
		Password: "secret",
		DB:       database.NewDatabase(&database.PostgresAdapter{}),
		Mutex:    &sync.Mutex{},
	}
	router(h).ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, `{ "status": "ok" }`, rec.Body.String())
}
