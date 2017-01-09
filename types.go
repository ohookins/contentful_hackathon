package main

// Content coming in from client to create a new entry
type clientEntry struct {
	Title      string `json:"title,omitempty"`
	URL        string `json:"url,omitempty"`
	Abstract   string `json:"abstract,omitempty"`
	ImageURL   string `json:"imageUrl,omitempty"`
	SourceName string `json:"sourceName,omitempty"`
	SourceURL  string `json:"sourceUrl,omitempty"`
}

// Much like the above, but formed in a way that the CMA expects.
type articleEntry struct {
	Fields articleFields `json:"fields"`
}

type articleFields struct {
	Title      nonLocalizedStringEntry `json:"title,omitempty"`
	URL        nonLocalizedStringEntry `json:"url,omitempty"`
	Abstract   nonLocalizedStringEntry `json:"abstract,omitempty"`
	ImageURL   nonLocalizedStringEntry `json:"imageUrl,omitempty"`
	SourceName nonLocalizedStringEntry `json:"sourceName,omitempty"`
	SourceURL  nonLocalizedStringEntry `json:"sourceUrl,omitempty"`
}

type nonLocalizedStringEntry struct {
	Data string `json:"en-US,omitempty"`
}

type createEntryResponse struct {
	Sys sys
}

type sys struct {
	ID      string
	Version uint
}

type deliveryResponse struct {
	Items []deliveryItem
}

type deliveryItem struct {
	Fields clientEntry
}
