package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	cf "github.com/Khaledgarbaya/contentful_go"
	"github.com/gorilla/mux"
)

const bindPort = 8080
const bodyTemplate = `
<html>
  <head><title>Test</title>
  <body>
    <ul>
{{ range $key, $value := . }}
      <li><strong>{{ $value.fields.title }}</strong>: {{ $value.fields.body }}</li>{{ end }}
    </ul>
  </body>
</html>
`

var (
	cdaToken = flag.String("cda-token", "", "Token or Access Key for the Content Delivery API")
	cmaToken = flag.String("cma-token", "", "Token or Access Key for the Content Management API")
	spaceID  = flag.String("space-id", "", "Space ID on Contentful")

	contentful cf.Contentful
	tmpl       *template.Template
)

// Content coming in from client to create a new entry
type clientEntry struct {
	Title string
	Body  string
}

func init() {
	flag.Parse()
	// TODO: Check that each flag is set
}

func main() {
	// Set up the Contentful client with the given space
	contentful = cf.New(*spaceID, *cdaToken)

	// Parse the template
	tmpl, _ = template.New("body").Parse(bodyTemplate)

	// Set up the routing
	r := mux.NewRouter()
	r.HandleFunc("/", summaryHandler)
	r.HandleFunc("/create", createHandler)
	http.Handle("/", r)

	// Set up the webserver
	listenPort := fmt.Sprintf(":%d", bindPort)
	log.Println("Now listening on", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read body of request", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	newEntry := clientEntry{}
	json.Unmarshal(body, &newEntry)

	http.Error(w, fmt.Sprintf("%v", newEntry), http.StatusOK)
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := contentful.GetEntries(map[string]string{})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, entries.Items)
}
