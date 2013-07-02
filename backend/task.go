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
	taskInstances := sess.ExistentCollection("TaskInstances")

	for _, id := range ids {
		task, err := tasks.Find(db.Cond{"tid": id})
		if err != nil {
			return
		}

		startTime := time.Now()
		endTime, err := time.Parse(timeFormat, task["duedate"].(string))

		events, _ := getEventsInRange(int(task["creator"].(int64)), startTime, endTime, sess)

		availability := NewSchedule(startTime, endTime)

		fmt.Println(availability.String())

		for _, event := range events {
			fmt.Println(event.Name, event.StartTime, event.StartTime.Add(event.Duration))
			availability.Subtract(event)
		}

		availability.FilterTimes(9, 17)
		availability.FilterMinimumTimeSpan(40 * time.Minute)
		times := availability.Sort()
		neededTime := time.Duration(task["timerequired"].(int64))

		for neededTime > 0 {
			scheduledTime := times[0]

			taskInstance := make(map[string]interface{})

			startTime := scheduledTime.start.Add(5 * time.Minute)
			length := 30 * time.Minute

			usedStart := true

			if 30 > neededTime {
				length = neededTime
			}

			if abs(12 - startTime.Hour()) > abs(12 - scheduledTime.end.Hour()) {
				startTime = scheduledTime.end.Add(-neededTime - 5)
				usedStart = false
			}

			taskInstance["when"] = startTime
			taskInstance["length"] = length
			taskInstance["tid"] = task["tid"]

			neededTime -= length

			taskInstances.Append(taskInstance)

			if usedStart {
				times[0].start.Add(length + 10 * time.Minute)
			} else {
				times[0].end.Add(-length - 10 * time.Minute)
			}

			times = append(times[1:], times[0])
		}

		fmt.Println(availability.String())
	}
}

func abs (a int) int {
	if a < 0 {
		return -a
	}
	return a
}
