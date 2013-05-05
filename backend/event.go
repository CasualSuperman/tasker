package main

import (
	"fmt"
	"github.com/gosexy/db"
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
	Name string            `json:"name"`
	StartTime time.Time    `json:"startTime"`
	Duration time.Duration `json:"duration"`
	Cid int				   `json:"cid"`

	allDay bool
	startDate time.Time
	endDate time.Time

	repeatType RepeatType
	repeatFrequency int

	repeatUntil time.Time
	days uint8
	fullWeek bool
}

func (e *Event) Parse(entry db.Item) {
	fmt.Printf("%+v\n", entry)

	fmt.Println(time.Parse(timeFormat, entry.GetString("start")))

	e.startDate, _ = time.Parse(timeFormat, entry.GetString("start"))
	e.endDate, _ = time.Parse(timeFormat,   entry.GetString("end"))

	e.repeatType = RepeatType(entry.GetInt("repeattype"))

	if e.repeatType != NoRepeat {
		e.repeatFrequency = int(entry.GetInt("repeatfrequency"))
		e.repeatUntil, _ = time.Parse(dateFormat, entry.GetString("repeatuntil"))
		e.days = uint8(entry.GetInt("days"))
		e.fullWeek = entry.GetBool("fullweek")
		e.repeatUntil = e.repeatUntil.AddDate(0, 0, 1) // The day after the last day we can be on.
	}

	e.Duration = time.Since(e.startDate) - time.Since(e.endDate)
	e.Name = entry.GetString("name")
	e.Cid = int(entry.GetInt("calendar"))
}

func (e *Event) FindInRange(start, end time.Time, resp chan Event) {
	switch e.repeatType {
		// If we don't repeat, make see if the original even occurs between the
		// start and end times. If it does, then send it. Always close the
		// channel.
		case NoRepeat:
			if e.startDate.Before(end) && e.endDate.After(start) {
				var eventCopy Event = *e
				eventCopy.StartTime = e.startDate
				resp <- eventCopy
			}

		// If it repeats every day, find the first time it happens in the range
		// and iterate through until we hit the end.
		case DailyRepeat:
			// Make sure the repeated date range is within the range we're scanning.
			if start.Before(e.repeatUntil) && end.After(e.startDate) {
				startDay := e.startDate
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
			if start.Before(e.repeatUntil) && end.After(e.startDate) {
				fmt.Println(e)
				startDate := e.startDate
				for startDate.Before(start) {
					startDate = startDate.AddDate(0, 0, e.repeatFrequency * 7)
				}
				// Made it to the first matching timespan.
				// Back it up because we went too far.
				startDate = startDate.AddDate(0, 0, e.repeatFrequency * -7)

				// This is an array of the number of days to add in a cycle
				// while hunting for hits.
				skips := makeSkips(e.days)
				skips[len(skips)-1] += 7 * (e.repeatFrequency-1)
				skipsIndex := 0

				for startDate.Before(start) || startDate.Before(e.startDate) {
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
	skips = skips[1:len(skips)-1]
	fmt.Println(skips)
	return skips
}
