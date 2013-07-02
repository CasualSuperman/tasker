package main

import (
	"fmt"
	"sort"
	"time"
)

type timeSpan struct {
	start, end time.Time
}

const tsFmt = "Jan 2 at 3:04pm"

func (t *timeSpan) String() string {
	return fmt.Sprintf("%s - %s", t.start.Format(tsFmt), t.end.Format(tsFmt))
}

func (t *timeSpan) Duration() time.Duration {
	return t.end.Sub(t.start)
}

func (s *schedule) String() string {
	r := ""
	for i, t := range *s {
		if i != 0 {
			r += ", "
		}
		r += t.String()
	}
	return r
}

type schedule []timeSpan

func NewSchedule(start, end time.Time) schedule {
	return schedule{timeSpan{start, end}}
}

func (s *schedule) Sort() []timeSpan {
	scp := append([]timeSpan(nil), (*s)...)
	cp := (*schedule)(&scp)
	sort.Sort(cp)

	return []timeSpan(*cp)
}

func (s *schedule) Len() int {
	return len(*s)
}

func (s schedule) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s schedule) Less(i, j int) bool {
	iDur := s[i].Duration()
	jDur := s[j].Duration()

	if iDur == jDur {
		return s[i].start.Before(s[j].start)
	} else {
		return iDur < jDur
	}
}

func (s *schedule) FilterTimes(startHour, endHour int) {
	firstDay := (*s)[0].start
	lastDay := (*s)[len(*s)-1].end

	currDay := firstDay.Truncate(time.Hour*24).UTC()
	currDayStart := currDay.Add(time.Duration(startHour)*time.Hour)
	currDayEnd := currDay.Add(time.Duration(endHour)*time.Hour)

	firstDayFromMidnight := Event{StartTime: currDay, Duration: time.Hour*time.Duration(startHour)}

	fmt.Println(firstDayFromMidnight)

	s.Subtract(firstDayFromMidnight)
	currDayStart = currDayStart.AddDate(0, 0, 1)

	for currDayEnd.Before(lastDay) {
		overNight := Event{StartTime: currDayEnd, Duration: currDayStart.Sub(currDayEnd)}

		fmt.Println(overNight)

		s.Subtract(overNight)

		currDayStart = currDayStart.AddDate(0, 0, 1)
		currDayEnd = currDayEnd.AddDate(0, 0, 1)
	}
}

func (s *schedule) Subtract(e Event) {
	fmt.Println("Before:", s.String())
	for i, section := range *s {
		if section.start.Before(e.StartTime) {
			fmt.Println("Section starts before Event")
			// We need to make sure the end of the e doesn't
			// intersect with the start of the section.
			// |---- Section ---?--
			//    ?--- Event ---?--
			if section.end.Before(e.StartTime) {
				fmt.Println("Section entirely before event")
				// If we're just after the section entirely, skip this
				// iteration.
				// |---- Section ----|
				//                     |--- Event ---?--
				continue
			} else {
				// |---- Section ----|
				//     |--- Event ---?--
				if section.end.After(e.StartTime.Add(e.Duration)) {
					fmt.Println("Event contained within section")
					// Here, the e is contained entirely within the
					// section.
					// |---- Section ----|
					//   |--- Event ---|
					// Turns into
					// |-|             |-|
					(*s) = append((*s)[:i],
						append([]timeSpan{
							{section.start, e.StartTime},
							{e.StartTime.Add(e.Duration), section.end}},
						(*s)[i+1:]...)...)
					continue
				} else {
					fmt.Println("Event tails section")
					// Otherwise, the e cuts off the remaining
					// time, and may cut into the next sections' times.
					// |---- Section ----|
					//     |--- Event ---?--
					// turns into
					// |---|
					(*s)[i].end = e.StartTime
					continue
				}
			}
		} else {
			fmt.Println("Event starts before section")
			// This means the e starts before the section starts.
			// We need to see if it ends before it begins as well.
			//         ?---- Section ----|
			// |--- Event ---?--
			if section.start.Before(e.StartTime.Add(e.Duration)) {
				fmt.Println("Event happens before section")
				// If it does, none of the remaining sections will
				// overlap, and we can continue on to the next e.
				//                 |---- Section ----|
				// |--- Event ---|
				continue
			} else if section.end.After(e.StartTime.Add(e.Duration)) {
				fmt.Println("Event heads section")
				// This means the e intersects the section in the
				// following manner:
				//     |---- Section ----|
				// |--- Event ---|
				// turns into
				//               |-------|
				(*s)[i].start = e.StartTime.Add(e.Duration)
				return
			} else {
				fmt.Println("Event contains section")
				// This means the e takes up the entire section,
				// and we should remove it.
				//   |---- Section ----|
				// |------- Event -------|
				(*s) = append((*s)[:i], (*s)[i+1:]...)
				continue
			}
		}
	}
}

func (s *schedule) FilterMinimumTimeSpan(d time.Duration) {
	for i := 0; i < len(*s); i++ {
		if (*s)[i].Duration() < d {
			*s = append((*s)[:i], (*s)[i+1:]...)
			i--
		}
	}
}
