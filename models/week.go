package models

import (
	"fmt"
	"sort"
	"time"
)

type TimeOfDay struct {
	Hour   int
	Minute int
}

func (t *TimeOfDay) Time() string {
	var h int
	var ampm string
	if t.Hour > 12 {
		h = t.Hour - 12
		ampm = "PM"
	} else {
		h = t.Hour
		ampm = "AM"
	}
	return fmt.Sprintf("%02d:%02d %s", h, t.Minute, ampm)
}

func (t *TimeOfDay) Time24H() string {
	return fmt.Sprintf("%02d%02d", t.Hour, t.Minute)
}

func (t *TimeOfDay) Difference(tod *TimeOfDay) int {
	diff := (t.Hour-tod.Hour)*60 + (t.Minute - tod.Minute)
	if diff < 0 {
		diff = -diff
	}

	return diff
}

func (t *TimeOfDay) LessThan(tod *TimeOfDay) bool {
	return t.Hour < tod.Hour || (t.Hour == tod.Hour && t.Minute < tod.Minute)
}

// ByAge implements sort.Interface for []*Activity based on
// the StartTime method.
type ByStartTime []*Activity

func (a ByStartTime) Len() int      { return len(a) }
func (a ByStartTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByStartTime) Less(i, j int) bool {
	iStart, _ := a[i].StartTime()
	jStart, _ := a[j].StartTime()
	iMinutes := iStart.Hour*60 + iStart.Minute
	jMinutes := jStart.Hour*60 + jStart.Minute
	return iMinutes < jMinutes
}

func eventsByWeekday(events []*Activity) (map[time.Weekday][]*Activity, error) {
	result := map[time.Weekday][]*Activity{}
	for i := time.Sunday; i <= time.Saturday; i++ {
		result[i] = []*Activity{}
	}

	for _, e := range events {
		weekday, err := e.Weekday()
		if err != nil {
			return result, err
		}

		result[weekday] = append(result[weekday], e)
	}

	for i := time.Sunday; i <= time.Saturday; i++ {
		sort.Sort(ByStartTime(result[i]))
	}

	return result, nil
}
