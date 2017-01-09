package main

// Content coming in from client to create a new entry
type clientEntry struct {
	Title    string `json:",omitempty"`
	URL      string `json:",omitempty"`
	Abstract string `json:",omitempty"`
}

// Much like the above, but formed in a way that the CMA expects.
type articleEntry struct {
	Fields articleFields `json:"fields"`
}

type articleFields struct {
	Title    nonLocalizedStringEntry `json:"title,omitempty"`
	URL      nonLocalizedStringEntry `json:"url,omitempty"`
	Abstract nonLocalizedStringEntry `json:"abstract,omitempty"`
}

type nonLocalizedStringEntry struct {
	Data string `json:"en-US,omitempty"`
}
