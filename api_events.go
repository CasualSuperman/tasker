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
				go events[i].FindInRange(
					time.Date(2013, time.April, 1, 0, 0, 0, 0, defaultLocation),
					time.Date(2013, time.May, 1, 0, 0, 0, 0, defaultLocation),
					eventChan)

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
			time.Date(2013, time.April, 1, 0, 0, 0, 0, defaultLocation),
			time.Date(2013, time.May, 1, 0, 0, 0, 0, defaultLocation),
		}
	}
	return &eventsList{
		[]Event{},
		time.Date(2013, time.April, 1, 0, 0, 0, 0, defaultLocation),
		time.Date(2013, time.May, 1, 0, 0, 0, 0, defaultLocation),
	}

}
