package main

import (
	"net/http"

	"github.com/JamesClonk/iRvisualizer/env"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/web"
)

func main() {
	port := env.Get("PORT", "8080")
	level := env.Get("LOG_LEVEL", "info")
	username := env.MustGet("AUTH_USERNAME")
	password := env.MustGet("AUTH_PASSWORD")

	log.Infoln("port:", port)
	log.Infoln("log level:", level)
	log.Infoln("auth username:", username)

	// start listener
	log.Fatalln(http.ListenAndServe(":"+port, web.NewRouter(username, password)))
}
