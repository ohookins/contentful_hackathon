package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	entryGetFmt = "https://cdn.contentful.com/spaces/%s/entries?content_type=article"
)

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := getEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Simplify the output format of all of the entries from the CDA.
	var outputEntries []clientEntry
	for _, value := range entries.Items {
		outputEntries = append(
			outputEntries,
			clientEntry{
				Title:    value.Fields.Title,
				URL:      value.Fields.URL,
				Abstract: value.Fields.Abstract,
			},
		)
	}

	// Prepare for writing to the output
	data, _ := json.Marshal(outputEntries)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getEntries() (deliveryResponse, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(entryGetFmt, spaceID), nil)
	req.Header.Add("Authorization", "Bearer "+cdaToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return deliveryResponse{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	entries := deliveryResponse{}
	_ = json.Unmarshal(body, &entries)

	return entries, nil
}
