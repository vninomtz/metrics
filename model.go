package main

import "time"

type View struct {
	Path    string    `json:"path"`
	IP      string    `json:"ip"`
	Agent   string    `json:"agent"`
	Created time.Time `json:"created"`
}
