package main

import (
	"log"
	"net/http"
)

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := contentful.GetEntries(map[string]string{})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, entries.Items)
}
