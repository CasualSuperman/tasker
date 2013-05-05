package main

import (
	"fmt"
	"github.com/gosexy/db"
	"net/http"
	"time"
)

const dateFormat = "2006-01-02"
const timeFormat = "2006-01-02 15:04"

var defaultLocation = time.UTC

func eventsInRange(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	session, _ := store.Get(req, "calendar")
	startDate, err := time.Parse(dateFormat, req.FormValue("start"))
	if err != nil {
		fmt.Println(startDate, err)
		return apiUserResponse{false, "Unable to parse start date.", http.StatusOK}
	}
	endDate, err := time.Parse(dateFormat, req.FormValue("end"))
	if err != nil {
		fmt.Println(endDate, err)
		return apiUserResponse{false, "Unable to parse end date.", http.StatusOK}
	}

	if val, ok := session.Values["logged-in"]; ok && val.(bool) {
		uid := int(session.Values["uid"].(int64))
		eventTable := sess.ExistentCollection("Events")

		eventResults, err := eventTable.FindAll(db.Cond{"creator": uid})

		events := make([]Event, len(eventResults))
		eventsInRange := make([]Event, 0)

		if err == nil {
			for i, event := range eventResults {
				eventChan := make(chan Event)
				(&events[i]).Parse(event)
				go events[i].FindInRange(startDate, endDate, eventChan)

				var ok bool = true
				var e Event

				for ok {
					e, ok = <-eventChan
					if ok {
						eventsInRange = append(eventsInRange, e)
					}
					fmt.Println(e)
				}
			}
		}
		return &eventsList{
			eventsInRange,
			startDate,
			endDate,
		}
	}
	return &eventsList{
		[]Event{},
		startDate,
		endDate,
	}

}
