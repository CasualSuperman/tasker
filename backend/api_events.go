package main

import (
	"fmt"
	"github.com/gosexy/db"
	"net/http"
	"time"
)

const dateFormat = "2006-01-02"
const timeFormat = "2006-01-02 15:04"
const formFormat = "2006-01-02 3:04pm"

var defaultLocation = time.UTC

func calendars(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	session, _ := store.Get(req, "calendar")

	if val, ok := session.Values["logged-in"]; ok && val.(bool) {
		uid := int(session.Values["uid"].(int64))
		calendarTable := sess.ExistentCollection("Calendars")

		calendarResults, err := calendarTable.FindAll(db.Cond{"owner": uid})

		calendars := calendarList(make([]calendar, len(calendarResults)))

		if err == nil {
			for i, cal := range calendarResults {
				calendars[i].Cid = int(cal.GetInt("cid"))
				calendars[i].Name = cal.GetString("name")
				calendars[i].Color = cal.GetString("color")
			}
		}
		return &calendars
	}
	return apiUserResponse{false, "Must be logged in.", http.StatusOK}
}

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

	if val, ok := session.Values["logged-in"]; !ok || !val.(bool) {
		return &eventsList{
			[]Event{},
			startDate,
			endDate,
		}
	}

	uid := int(session.Values["uid"].(int64))
	eventTable := sess.ExistentCollection("Events")

	eventResults, err := eventTable.FindAll(db.Cond{"creator": uid})

	events := make([]Event, len(eventResults))
	eventsInRange := make([]Event, 0)

	if err == nil {
		for i, event := range eventResults {
			eventChan := make(chan Event)
			(&events[i]).ParseDB(event)
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

func createEvent(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	session, _ := store.Get(req, "calendar")

	if val, ok := session.Values["logged-in"]; !ok || !val.(bool) {
		return apiUserResponse{
			false,
			"Please register to create an event.",
			http.StatusOK,
		}
	}

	uid := int(session.Values["uid"].(int64))
	event, errFields, errMsgs := ParseHTTP(req)
	event["creator"] = uid
	event["calendar"] = 1

	fmt.Println(event)

	if len(errFields) > 0 {
		_ = errMsgs
	} else {
		eventTable := sess.ExistentCollection("Events")
		_, err := eventTable.Append(event)
		if err != nil {
			return apiUserResponse{false, err.Error(), http.StatusOK}
		}
	}

	return apiUserResponse{}
}
