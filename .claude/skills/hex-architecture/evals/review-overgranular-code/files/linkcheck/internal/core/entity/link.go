package entity

import "time"

type Link struct {
	URL        string    `json:"url"`
	StatusCode int       `json:"status_code"`
	CheckedAt  time.Time `json:"checked_at"`
}

func (l Link) IsBroken() bool {
	return l.StatusCode >= 400
}

// State is the resume checkpoint of a check run.
type State struct {
	CurrentFile   string // name of the sitemap file being checked
	CheckedURLs   []string
}
