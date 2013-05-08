package main

import (
	"fmt"
	"github.com/gosexy/db"
	"net/http"
	"strconv"
	"time"
)

type RepeatType uint8

const (
	NoRepeat RepeatType = iota
	DailyRepeat
	WeeklyRepeat
	MonthlyRepeat
	YearlyRepeat
)

type Event struct {
	Name      string        `json:"name"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
	Calendar  int           `json:"cid"`
	Eid       int           `json:"eid"`
	creator   int

	allDay bool

	start time.Time
	end   time.Time

	repeatType      RepeatType
	repeatFrequency int

	repeatUntil time.Time
	days        uint8
	fullWeek    bool
	weekOfMonth int
}

func (e *Event) ParseDB(entry db.Item) {
	e.start, _ = time.Parse(timeFormat, entry.GetString("start"))
	e.end, _ = time.Parse(timeFormat, entry.GetString("end"))

	e.repeatType = RepeatType(entry.GetInt("repeattype"))

	if e.repeatType != NoRepeat {
		e.repeatFrequency = int(entry.GetInt("repeatfrequency"))
		e.repeatUntil, _ = time.Parse(dateFormat, entry.GetString("repeatuntil"))
		e.days = uint8(entry.GetInt("days"))
		e.fullWeek = entry.GetBool("fullweek")
		e.repeatUntil = e.repeatUntil.AddDate(0, 0, 1) // The day after the last day we can be on.
	}

	e.Duration = time.Since(e.start) - time.Since(e.end)
	e.Name = entry.GetString("name")
	e.Calendar = int(entry.GetInt("calendar"))
	e.Eid = int(entry.GetInt("eid"))
}

func ParseHTTP(req *http.Request) (map[string]interface{}, []string, []string) {
	e := make(map[string]interface{})
	errFields := make([]string, 0)
	errErrors := make([]string, 0)

	// Get the event name
	e["name"] = req.FormValue("name")
	if req.FormValue("name") == "" {
		errFields = append(errFields, "name")
		errErrors = append(errErrors, "Name is required.")
	}

	var err error

	// Get the calendar we want to use.
	e["calendar"], err = getIntFromHTTP(req, "calendar")
	if err != nil {
		errFields = append(errFields, "calendar")
		errErrors = append(errErrors, "Invalid calendar.")
	}

	// Get the event start time
	startStr := req.FormValue("startDate_submit") + " " + req.FormValue("startTime")
	startTime, err := time.Parse(formFormat, startStr)
	e["start"] = startTime.Format(timeFormat)
	if err != nil {
		fmt.Println(err)
		errFields = append(errFields, "startTime")
		errErrors = append(errErrors, "Start date format incorrect.")
	}

	// Get the event end time
	endStr := req.FormValue("endDate_submit") + " " + req.FormValue("endTime")
	endTime, err := time.Parse(formFormat, endStr)
	e["end"] = endTime.Format(timeFormat)
	if err != nil {
		errFields = append(errFields, "endTime")
		errErrors = append(errErrors, "End date format incorrect.")
	}

	// Get the repeat type
	repeatStr := req.FormValue("frequency")
	switch repeatStr {
	case "none":
		e["repeatType"] = NoRepeat
	case "daily":
		e["repeatType"] = DailyRepeat
	case "weekly":
		e["repeatType"] = WeeklyRepeat
	case "monthly":
		e["repeatType"] = MonthlyRepeat
	case "yearly":
		e["repeatType"] = YearlyRepeat
	default:
		errFields = append(errFields, "frequency")
		errErrors = append(errErrors, "Unrecognized repeat frequency.")
	}

	e["allDay"] = req.FormValue("allDay") == "on"

	// We only need to get these if we're gonna repeat.
	if e["repeatType"] != NoRepeat {
		// The number of things to skip.
		e["repeatFrequency"], err = getIntFromHTTP(req, "skip")
		if err != nil {
			errFields = append(errFields, "skip")
			errErrors = append(errErrors, "Improper skip amount.")
		}

		doesEnd := req.FormValue("ends")

		switch doesEnd {
		case "never":
			e["repeatUntil"] = nil
		case "afterN":
			errFields = append(errFields, "afterN")
			errErrors = append(errErrors, "After X times unavailable.")
		case "afterDate":
			lastStr := req.FormValue("afterDate_submit")
			lastTime, err := time.Parse(dateFormat, lastStr)
			e["repeatUntil"] = lastTime.Format(dateFormat)
			if err != nil {
				errFields = append(errFields, "afterDate")
				errErrors = append(errErrors, "Last date format incorrect.")
			}
		default:
			errFields = append(errFields, "ends")
			errErrors = append(errErrors, "Unrecognized repeat stop.")
		}

		if (e["repeatType"] == MonthlyRepeat && req.FormValue("repeatByMonth") == "day") ||
			e["repeatType"] == WeeklyRepeat {
			// Figure out the bitmask of which days we repeat on.
			days := req.Form["daysOfWeek"]
			dayBitMask := 0
			daySlice := []string{"Su", "M", "Tu", "W", "Th", "F", "Sa"}
			for _, str := range daySlice {
				dayBitMask <<= 1
				for _, matchStr := range days {
					if str == matchStr {
						dayBitMask |= 1
					}
				}
			}
			e["days"] = uint8(dayBitMask)

			if dayBitMask == 0 {
				errFields = append(errFields, "daysOfWeek")
				errErrors = append(errErrors, "Please select at least one day of the week.")
			}

			if e["repeatType"] == MonthlyRepeat {
				e["fullWeek"] = req.FormValue("fullWeek") == "on"

				// Figure out which week of the month we repeat on.
				weekInMonthStr := req.FormValue("weekInMonth")
				weekInMonthInt, err := strconv.Atoi(weekInMonthStr)
				if err != nil {
					errFields = append(errFields, "weekInMonth")
					errErrors = append(errErrors, "Somehow you messed up the form.")
				}
				e["weekOfMonth"] = weekInMonthInt
			}
		}
	}

	return e, errFields, errErrors
}

func (e *Event) FindInRange(start, end time.Time, resp chan Event) {
	switch e.repeatType {
	// If we don't repeat, make see if the original even occurs between the
	// start and end times. If it does, then send it. Always close the
	// channel.
	case NoRepeat:
		if e.start.Before(end) && e.end.After(start) {
			var eventCopy Event = *e
			eventCopy.StartTime = e.start
			resp <- eventCopy
		}

	// If it repeats every day, find the first time it happens in the range
	// and iterate through until we hit the end.
	case DailyRepeat:
		// Make sure the repeated date range is within the range we're scanning.
		if start.Before(e.repeatUntil) && end.After(e.start) {
			startDay := e.start
			for startDay.Before(start) {
				startDay = startDay.AddDate(0, 0, e.repeatFrequency)
			}
			for startDay.Before(end) {
				var eventInstance Event = *e
				eventInstance.StartTime = startDay
				resp <- eventInstance
				startDay = startDay.AddDate(0, 0, e.repeatFrequency)
			}
		}

	case WeeklyRepeat:
		if start.Before(e.repeatUntil) && end.After(e.start) {
			fmt.Println(e)
			startDate := e.start
			for startDate.Before(start) {
				startDate = startDate.AddDate(0, 0, e.repeatFrequency*7)
			}
			// Made it to the first matching timespan.
			// Back it up because we went too far.
			startDate = startDate.AddDate(0, 0, e.repeatFrequency*-7)

			// This is an array of the number of days to add in a cycle
			// while hunting for hits.
			skips := makeSkips(e.days)
			skips[len(skips)-1] += 7 * (e.repeatFrequency - 1)
			skipsIndex := 0

			for startDate.Before(start) || startDate.Before(e.start) {
				startDate = startDate.AddDate(0, 0, skips[skipsIndex])

				skipsIndex++
				if skipsIndex == len(skips) {
					skipsIndex = 0
				}
			}

			// Found the first potential match.

			for startDate.Before(end) && startDate.Before(e.repeatUntil) {
				// Sending a match.
				eventInstance := *e
				eventInstance.StartTime = startDate

				resp <- eventInstance

				// Adding for the next match.

				startDate = startDate.AddDate(0, 0, skips[skipsIndex])

				skipsIndex++
				if skipsIndex == len(skips) {
					skipsIndex = 0
				}
			}
		}
	}
	close(resp)
}

func makeSkips(bits uint8) []int {
	bits <<= 1
	bits |= 1
	skips := make([]int, 1)
	for bits != 0 {
		if (bits & 0x80) == 0x80 {
			skips = append(skips, 1)
		} else {
			skips[len(skips)-1]++
		}
		bits <<= 1
	}
	skips[len(skips)-2] += skips[0]
	skips = skips[1 : len(skips)-1]
	fmt.Println(skips)
	return skips
}

func getIntFromHTTP(req *http.Request, field string) (int, error) {
	tempStr := req.FormValue(field)
	return strconv.Atoi(tempStr)
}
