package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	cf "github.com/Khaledgarbaya/contentful_go"
	"github.com/gorilla/mux"
	"github.com/streadway/handy/report"
)

const bodyTemplate = `
<html>
  <head><title>Test</title>
  <body>
    <ul>
{{ range $key, $value := . }}
      <li><strong>{{ $value.fields.title }}</strong>: {{ $value.fields.url }}</li>{{ end }}
    </ul>
  </body>
</html>
`

var (
	cdaToken = os.Getenv("CFHACK_CDA_TOKEN")
	cmaToken = os.Getenv("CFHACK_CMA_TOKEN")
	spaceID  = os.Getenv("CFHACK_SPACE_ID")
	bindPort = os.Getenv("PORT")

	contentful cf.Contentful
	tmpl       *template.Template
)

func main() {
	// Set up the Contentful client with the given space
	contentful = cf.New(spaceID, cdaToken)

	// Parse the template
	tmpl, _ = template.New("body").Parse(bodyTemplate)

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
