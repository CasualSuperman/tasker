package main

import (
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
			task["timeRequired"] = hours * 60 + mins
		}
	}

	dueStr := req.FormValue("dueDate_submit") + " " + req.FormValue("dueTime")
	dueDate, err := time.Parse(formFormat, dueStr)
	task["dueDate"] = dueDate

	if err != nil {
		errFields = append(errFields, "dueDate")
		errMsgs = append(errMsgs, "Due date format incorrect.")
	}

	if len(errFields) > 0 {
		return &apiFormResponse{false, errFields, errMsgs}
	}

	taskTable := sess.ExistentCollection("Tasks")
	_, err = taskTable.Append(task)
	if err != nil {
		println(err.Error())
		return &apiFormResponse{false, nil, nil}
	}
	return &apiFormResponse{true, nil, nil}
}
