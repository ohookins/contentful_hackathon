package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/streadway/handy/report"
)

var (
	cdaToken = os.Getenv("CFHACK_CDA_TOKEN")
	cmaToken = os.Getenv("CFHACK_CMA_TOKEN")
	spaceID  = os.Getenv("CFHACK_SPACE_ID")
	bindPort = os.Getenv("PORT")
)

func main() {
	// Set up the routing
	r := mux.NewRouter()
	r.HandleFunc("/", summaryHandler)
	r.HandleFunc("/create", createHandler)
	http.Handle("/", report.JSON(os.Stdout, r))

	// Set up the webserver
	listenPort := fmt.Sprintf(":%s", bindPort)
	log.Println("Now listening on", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
