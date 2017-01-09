package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	entryCreationFmt = "https://api.contentful.com/spaces/%s/entries"
	entryPublishFmt  = "https://api.contentful.com/spaces/%s/entries/%s/published"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read body of request", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// Transform the entry received from the client into something the CMA can
	// handle
	buf := transformEntry(body)

	// Create the new entry on the CMA
	entryID, version, err := createEntry(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Now publish it
	err = publishEntry(entryID, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, fmt.Sprintf("created and published entry %s", entryID), http.StatusCreated)
}

func transformEntry(incomingEntry []byte) *bytes.Buffer {
	newEntry := clientEntry{}
	json.Unmarshal(incomingEntry, &newEntry)
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
	return buf
}

func createEntry(buf *bytes.Buffer) (string, uint, error) {
	req, _ := http.NewRequest("POST", fmt.Sprintf(entryCreationFmt, spaceID), buf)
	req.Header.Add("Authorization", "Bearer "+cmaToken)
	req.Header.Add("Content-Type", "application/vnd.contentful.management.v1+json")
	req.Header.Add("X-Contentful-Content-Type", "article")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	// Buffer the entire body for reading the entry ID
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("%v", body)
		return "", 0, fmt.Errorf("unexpected return code from CMA: %d", resp.StatusCode)
	}

	entryID, version := getCreatedEntryIDAndVersion(body)
	return entryID, version, nil
}

func getCreatedEntryIDAndVersion(body []byte) (string, uint) {
	createdEntry := createEntryResponse{}
	json.Unmarshal(body, &createdEntry)

	return createdEntry.Sys.ID, createdEntry.Sys.Version
}

func publishEntry(entryID string, version uint) error {
	req, _ := http.NewRequest("PUT", fmt.Sprintf(entryPublishFmt, spaceID, entryID), nil)
	req.Header.Add("Authorization", "Bearer "+cmaToken)
	req.Header.Add("Content-Type", "application/vnd.contentful.management.v1+json")
	req.Header.Add("X-Contentful-Version", fmt.Sprintf("%d", version))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Printf("%s\n", body)
	return nil
}
