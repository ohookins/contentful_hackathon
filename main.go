package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	cf "github.com/Khaledgarbaya/contentful_go"
	"github.com/gorilla/mux"
)

const (
	entryCreationFmt = "https://api.contentful.com/spaces/%s/entries"
	bodyTemplate     = `
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
)

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
	http.Handle("/", r)

	// Set up the webserver
	listenPort := fmt.Sprintf(":%s", bindPort)
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

	// Transform the entry received from the client into something the CMA can
	// handle
	newEntry := clientEntry{}
	json.Unmarshal(body, &newEntry)
	entry := &articleEntry{
		Fields: articleFields{
			Title:    nonLocalizedStringEntry{newEntry.Title},
			URL:      nonLocalizedStringEntry{newEntry.URL},
			Abstract: nonLocalizedStringEntry{newEntry.Abstract},
		},
	}
	jsonEntry, _ := json.Marshal(entry)
	buf := bytes.NewBuffer(jsonEntry)

	log.Printf("%s\n", buf)

	// Create the new entry on the CMA
	req, _ := http.NewRequest("POST", fmt.Sprintf(entryCreationFmt, spaceID), buf)
	req.Header.Add("Authorization", "Bearer "+cmaToken)
	req.Header.Add("Content-Type", "application/vnd.contentful.management.v1+json")
	req.Header.Add("X-Contentful-Content-Type", "article")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("%s\n", body)

		http.Error(
			w,
			fmt.Sprintf("unexpected return code from CMA: %d", resp.StatusCode),
			http.StatusInternalServerError,
		)
		return
	}

	http.Error(w, fmt.Sprintf("created entry for %s", newEntry.URL), http.StatusCreated)
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
