package main

import (
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
	Eid       int           `json:"eid"`
	Name      string        `json:"name"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
	Calendar  int           `json:"cid"`
	creator   int

	AllDay bool `json:"allDay"`

	start time.Time
	end   time.Time

	repeatType      RepeatType
	repeatFrequency int

	repeatUntil *time.Time
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
		repeatUntil, err := time.Parse(dateFormat, entry.GetString("repeatuntil"))
		if err == nil {
			repeatUntil = repeatUntil.AddDate(0, 0, 1) // The day after the last day we can be on.
			e.repeatUntil = &repeatUntil
		}
		e.days = uint8(entry.GetInt("days"))
		e.fullWeek = entry.GetBool("fullweek")
	}

	e.Duration = time.Since(e.start) - time.Since(e.end)
	e.Name = entry.GetString("name")
	e.Calendar = int(entry.GetInt("calendar"))
	e.Eid = int(entry.GetInt("eid"))
	e.AllDay = entry.GetBool("allday")
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

// TODO: Implement MonthlyRepeat
// TODO: Implement Repeat Counts instead of dates or forevers
func (e *Event) FindInRange(start, end time.Time, resp chan Event) {
	defer close(resp)

	if e.repeatType != NoRepeat {
		if e.start.After(end) {
			return
		}
		if e.repeatUntil != nil && e.repeatUntil.Before(start) {
			return
		}
	}

	// If the event can't happen in the range, we can just quit here.
	switch e.repeatType {

	case NoRepeat:
		// If we don't repeat, make see if the original even occurs between the
		// start and end times. If it does, then send it. Always close the
		// channel.
		if e.start.After(end) || e.end.Before(start) {
			return
		}
		var eventCopy Event = *e
		eventCopy.StartTime = e.start
		resp <- eventCopy

	case DailyRepeat:
		// If it repeats every day, find the first time it happens in the range
		// and iterate through until we hit the end.
		startDay := e.start
		for startDay.Before(start) {
			startDay = startDay.AddDate(0, 0, e.repeatFrequency)
		}

		lastDay := end
		if e.repeatUntil != nil {
			lastDay = *e.repeatUntil
		}

		for startDay.Before(end) && startDay.Before(lastDay) {
			var eventInstance Event = *e
			eventInstance.StartTime = startDay
			resp <- eventInstance
			startDay = startDay.AddDate(0, 0, e.repeatFrequency)
		}

	case WeeklyRepeat:
		// Find the first week that the event happens in within our range.
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

		// While the event is before when the event starts, or when the range
		// starts, go to the next one.
		for startDate.Before(start) || startDate.Before(e.start) {
			startDate = startDate.AddDate(0, 0, skips[skipsIndex])

			skipsIndex++
			if skipsIndex == len(skips) {
				skipsIndex = 0
			}
		}

		// Found the first potential match.

		// Now find all the matches and send them down the channel.
		for startDate.Before(end) && (e.repeatUntil == nil || startDate.Before(*e.repeatUntil)) {
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
	case YearlyRepeat:
		dateThisYear := e.start
		years := start.Year() - dateThisYear.Year()

		// If the starting point of the event is before the year the span
		// contains.
		for years < 0 {
			years += e.repeatFrequency
		}

		// Move the event to this year or after it
		dateThisYear = dateThisYear.AddDate(years, 0, 0)

		// We're now after start, make sure we're before end.
		for dateThisYear.Before(end) && dateThisYear.Add(e.Duration).After(start) {
			// Sending a match
			eventThisYear := *e
			eventThisYear.StartTime = dateThisYear

			resp <- eventThisYear

			dateThisYear = dateThisYear.AddDate(e.repeatFrequency, 0, 0)
		}
	}
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
	return skips
}

func getIntFromHTTP(req *http.Request, field string) (int, error) {
	tempStr := req.FormValue(field)
	return strconv.Atoi(tempStr)
}

func getEventsInRange(uid int, startDate, endDate time.Time, sess db.Database) ([]Event, error) {
	eventTable := sess.ExistentCollection("Events")

	eventResults, err := eventTable.FindAll(db.Cond{"creator": uid})

	events := make([]Event, len(eventResults))
	var eventsInRange []Event
	eventChan := make(chan Event, 10)
	doneChan := make(chan bool)

	if err != nil {
		return nil, err
	}
	waitingFor := len(eventResults)

	for i, event := range eventResults {
		(&events[i]).ParseDB(event)
		personalEventChan := make(chan Event)

		// Start finding events but continue looping.
		go events[i].FindInRange(startDate, endDate, personalEventChan)

		// Watch for the returned events. If the channel gets closed, tell
		// the waitGroup we're done. Otherwise, keep watching.
		go func() {
			var e Event
			ok := true

			for ok {
				if e, ok = <-personalEventChan; ok {
					eventChan <- e
				}
			}
			doneChan <- true
		}()
	}
	for waitingFor > 0 {
		select {
		case <-doneChan:
			waitingFor--
		case e := <-eventChan:
			eventsInRange = append(eventsInRange, e)
		}
	}
	return eventsInRange, nil
}
