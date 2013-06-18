package main

import (
	"github.com/gosexy/db"
	"net/http"
	"strconv"
	"time"
)

const dateFormat = "2006-01-02"
const timeFormat = "2006-01-02 15:04"
const formFormat = "2006-01-02 3:04 PM"

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
		return apiUserResponse{false, "Unable to parse start date.", http.StatusOK}
	}
	endDate, err := time.Parse(dateFormat, req.FormValue("end"))
	if err != nil {
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
			}
		}
	}

	return &eventsList{
		eventsInRange,
		startDate,
		endDate,
	}
}

func dumpEvent(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	session, _ := store.Get(req, "calendar")

	if val, ok := session.Values["logged-in"]; !ok || !val.(bool) {
		return apiUserResponse{
			false,
			"Please login to get event info.",
			http.StatusOK,
		}
	}

	uid := int(session.Values["uid"].(int64))
	eid, err := strconv.Atoi(req.FormValue("eid"))
	if err != nil {
		return apiUserResponse{
			false,
			"Unable to parse event id.",
			http.StatusOK,
		}
	}

	events := sess.ExistentCollection("Events")
	event, err := events.Find(db.Cond{"eid": eid}, db.Cond{"creator": uid})

	if err == nil {
		return apiRawResponse{
			true,
			event,
		}
	}
	return apiUserResponse{
		true,
		"This event doesn't exist or you don't have permission to edit it.",
		http.StatusOK,
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

	if len(errFields) > 0 {
		return &apiFormResponse{false, errFields, errMsgs}
	} else {
		if checkUserOwnsCalendar(sess, uid, event["calendar"].(int)) {
			eventTable := sess.ExistentCollection("Events")
			_, err := eventTable.Append(event)
			if err != nil {
				return &apiFormResponse{false, nil, nil}
			}
			return &apiFormResponse{true, nil, nil}
		} else {
			return &apiFormResponse{
				false,
				[]string{"calendar"},
				[]string{"You don't have permission to use that calendar."},
			}
		}
	}

	return apiUserResponse{}
}

func updateEvent(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
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
	eid, err := strconv.Atoi(req.FormValue("eid"))
	event["eid"] = eid

	if len(errFields) > 0 || err != nil {
		return &apiFormResponse{false, errFields, errMsgs}
	} else {
		if checkUserOwnsCalendar(sess, uid, event["calendar"].(int)) &&
			checkUserOwnsEvent(sess, uid, event["eid"].(int)) {
			eventTable := sess.ExistentCollection("Events")
			err := eventTable.Update(db.Cond{"eid": eid}, db.Set(event))
			if err != nil {
				return &apiFormResponse{false, nil, nil}
			}
			return &apiFormResponse{true, nil, nil}
		} else {
			return &apiFormResponse{
				false,
				[]string{"calendar"},
				[]string{"You don't have permission to use that calendar."},
			}
		}
	}

	return apiUserResponse{}
}

func checkUserOwnsCalendar(sess db.Database, uid, cal int) bool {
	calendars := sess.ExistentCollection("Calendars")

	num, err := calendars.Count(db.Cond{"cid": cal}, db.Cond{"owner": uid})
	if err != nil {
		return false
	} else if num > 0 {
		return true
	}

	sharedCalendars := sess.ExistentCollection("CalendarShares")
	num, err = sharedCalendars.Count(db.Cond{"cid": cal}, db.Cond{"uid": uid})
	if err != nil {
		return false
	} else if num > 0 {
		return true
	}

	return false
}

func checkUserOwnsEvent(sess db.Database, uid, eid int) bool {
	events := sess.ExistentCollection("Events")

	num, err := events.Count(db.Cond{"eid": eid}, db.Cond{"creator": uid})
	if err != nil {
		return false
	} else if num > 0 {
		return true
	}

	sharedEvents := sess.ExistentCollection("EventShares")
	num, err = sharedEvents.Count(db.Cond{"eid": eid}, db.Cond{"uid": uid})
	if err != nil {
		return false
	} else if num > 0 {
		return true
	}

	return false
}
