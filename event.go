package main

import (
	"time"
)

type Event struct {
	Name string            `json:"name"`
	StartTime time.Time    `json:"startTime"`
	Duration time.Duration `json:"duration"`
}
