package main

import (
	"github.com/gosexy/db"
	"net/http"
	"time"
)

var defaultLocation, _ = time.LoadLocation("America/New_York")

func eventsInRange(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	return &eventsList{
		[]Event{
			Event{
				"Capstone Meeting",
				time.Date(2013, time.April, 16, 10, 0, 0, 0, defaultLocation),
				time.Duration(15 * time.Minute),
			},
		},
		time.Date(2013, time.April, 1, 0, 0, 0, 0, defaultLocation),
		time.Date(2013, time.May, 1, 0, 0, 0, 0, defaultLocation),
	}
}
