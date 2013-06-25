package main

import (
"fmt"
	"github.com/gosexy/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Tid          int           `json:"tid"`
	Name         string        `json:"name"`
	TimeRequired time.Duration `json:"totalTime"`
	TimeInvested time.Duration `json:"givenTime"`
	Due          time.Time     `json:"due"`
}

func createTask(res http.ResponseWriter, req *http.Request, sess db.Database) apiResponse {
	session, _ := store.Get(req, "calendar")

	if val, ok := session.Values["logged-in"]; !ok || !val.(bool) {
		return apiUserResponse{
			false,
			"Please register to create an event.",
			http.StatusOK,
		}
	}

	uid := int(session.Values["uid"].(int64))
	task := make(map[string]interface{})
	var errFields, errMsgs []string
	task["creator"] = uid
	task["name"] = req.FormValue("name")
	task["timeInvested"] = 0

	timeRequiredStrs := strings.Split(req.FormValue("duration"), ":")

	if len(timeRequiredStrs) != 2 {
		errFields = append(errFields, "duration")
		errMsgs = append(errMsgs, "Please pick a proper duration.")
	} else {
		hoursStr, minsStr := timeRequiredStrs[0], timeRequiredStrs[1]
		hours, hErr := strconv.Atoi(hoursStr)
		mins, mErr := strconv.Atoi(minsStr)

		if hErr != nil || mErr != nil {
			errFields = append(errFields, "duration")
			errMsgs = append(errMsgs, "Please pick a proper duration.")
		} else {
			task["timeRequired"] = hours*60 + mins
		}
	}

	dueStr := req.FormValue("dueDate_submit") + " " + req.FormValue("dueTime")
	dueDate, err := time.Parse(formFormat, dueStr)
	task["dueDate"] = dueDate.Format(timeFormat)

	if err != nil {
		errFields = append(errFields, "dueDate")
		errMsgs = append(errMsgs, "Due date format incorrect.")
	}

	if len(errFields) > 0 {
		return &apiFormResponse{false, errFields, errMsgs}
	}

	taskTable := sess.ExistentCollection("Tasks")
	id, err := taskTable.Append(task)
	if err != nil {
		println(err.Error())
		return &apiFormResponse{false, nil, nil}
	}
	createTaskInstances(id, sess)
	return &apiFormResponse{true, nil, nil}
}

func createTaskInstances(ids []db.Id, sess db.Database) {
	tasks := sess.ExistentCollection("Tasks")

	for _, id := range ids {
		task, err := tasks.Find(db.Cond{"tid": id})
		if err != nil {
			return
		}

		startTime := time.Now()
		fmt.Println(startTime)
		endTime, err := time.Parse(timeFormat, task["duedate"].(string))

		events, _ := getEventsInRange(int(task["creator"].(int64)), startTime, endTime, sess)

		availability := [][2]time.Time{{startTime, endTime}}

		fmt.Println(availability)

		taskEvents:
		for _, event := range events {
			fmt.Println(event.Name, event.StartTime, event.StartTime.Add(event.Duration))
			eventSections:
			for i, section := range availability {
				if section[0].Before(event.StartTime) {
					// We need to make sure the end of the event doesn't
					// intersect with the start of the section.
					// |---- Section ---?--
					//    ?--- Event ---?--
					if section[1].Before(event.StartTime) {
						// If we're just after the section entirely, skip this
						// iteration.
						// |---- Section ----|
						//                     |--- Event ---?--
						continue eventSections
					} else {
						// |---- Section ----|
						//     |--- Event ---?--
						if section[1].After(event.StartTime.Add(event.Duration)) {
							// Here, the event is contained entirely within the
							// section.
							// |---- Section ----|
							//   |--- Event ---|
							// Turns into
							// |-|             |-|
							availability = append(availability[:i], append([][2]time.Time{{section[0], event.StartTime}, {event.StartTime.Add(event.Duration), section[1]}}, availability[i+1:]...)...)
							continue taskEvents
						} else {
							// Otherwise, the event cuts off the remaining
							// time, and may cut into the next sections' times.
							// |---- Section ----|
							//     |--- Event ---?--
							// turns into
							// |---|
							availability[i][1] = event.StartTime
							continue eventSections
						}
					}
				} else {
					// This means the event starts before the section starts.
					// We need to see if it ends before it begins as well.
					//         ?---- Section ----|
					// |--- Event ---?--
					if section[0].Before(event.StartTime.Add(event.Duration)) {
						// If it does, none of the remaining sections will
						// overlap, and we can continue on to the next event.
						//                 |---- Section ----|
						// |--- Event ---|
						continue taskEvents
					} else if section[1].After(event.StartTime.Add(event.Duration)) {
						// This means the event intersects the section in the
						// following manner:
						//     |---- Section ----|
						// |--- Event ---|
						// turns into
						//               |-------|
						availability[i][0] = event.StartTime.Add(event.Duration)
						continue taskEvents
					} else {
						// This means the event takes up the entire section,
						// and we should remove it.
						//   |---- Section ----|
						// |------- Event -------|
						availability = append(availability[:i], availability[i+1:]...)
						continue eventSections
					}
				}
			}
		}

		fmt.Println(availability)
	}
}
