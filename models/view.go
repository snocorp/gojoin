package models

import (
	"fmt"
	"strings"
	"time"
)

type ViewEvent struct {
	Activity *Activity

	StartTime string
	Duration  int
	Offset    int
	Span      int
	BgColor   string

	prevEvent *ViewEvent
}

func (ve *ViewEvent) String() string {
	return fmt.Sprint(ve.Activity.Id)
}

type WeekDay struct {
	Name      string
	ShortName string
}

type Time struct {
	Name    string
	Code    string
	GridRow int
}

type WeekdayView struct {
	Span        int
	GridColumns []int
	Events      []*ViewEvent
}

type CenterView struct {
	CenterId   string
	CenterName string

	GridColumns string

	Weekdays []*WeekdayView
}

type View struct {
	Centers []*CenterView

	Days []WeekDay

	Times []Time
}

func weekdays() []WeekDay {
	return []WeekDay{
		{"Sunday", "Sun"},
		{"Monday", "Mon"},
		{"Tuesday", "Tue"},
		{"Wednesday", "Wed"},
		{"Thursday", "Thu"},
		{"Friday", "Fri"},
		{"Saturday", "Sat"},
	}
}

func times() []Time {
	list := []Time{}
	gridRow := 2
	for h := 9; h < 21; h++ {
		for m := 0; m < 60; m += 15 {
			t := &TimeOfDay{h, m}
			list = append(list, Time{Name: t.Time(), Code: t.Time24H(), GridRow: gridRow})
			gridRow++
		}
	}

	return list
}

/*

A
AB
 BCD
  CD

*/

func NewViewEvent(a *Activity, prevEvent *ViewEvent, colourMap map[string]string) (*ViewEvent, error) {
	offset := 0
	if prevEvent != nil {
		overlaps, err := prevEvent.Activity.Overlaps(a)
		if err != nil {
			return nil, err
		}

		if overlaps {
			overlapMap := map[int]*ViewEvent{
				prevEvent.Offset: prevEvent,
			}
			offsetEvent := prevEvent
			for offsetEvent.prevEvent != nil {
				prevOverlaps, err := offsetEvent.prevEvent.Activity.Overlaps(a)
				if err != nil {
					return nil, err
				}

				if prevOverlaps {
					overlapMap[offsetEvent.prevEvent.Offset] = offsetEvent.prevEvent
				}

				offsetEvent = offsetEvent.prevEvent
			}

			var i int
			ok := true
			for i = 0; ok; i++ {
				_, ok = overlapMap[i]
			}

			offset = i - 1
		} else {
			prevEvent = nil
		}
	}

	st, err := a.StartTime()
	if err != nil {
		return nil, err
	}

	et, err := a.EndTime()
	if err != nil {
		return nil, err
	}

	return &ViewEvent{
		Activity:  a,
		StartTime: st.Time24H(),
		Duration:  et.Difference(&st),
		Offset:    offset,
		Span:      1,
		BgColor:   colourMap[a.Name],
		prevEvent: prevEvent,
	}, nil
}

func NewCenterView(cw *CenterWeek, days []WeekDay) (*CenterView, error) {
	colours := []string{
		"rgb(234, 153, 153)",
		"rgb(249, 203, 156)",
		"rgb(255, 229, 153)",
		"rgb(182, 215, 168)",
		"rgb(162, 196, 201)",
		"rgb(164, 194, 244)",
	}
	colourMap := map[string]string{}
	colorIndex := 0
	for _, e := range cw.Events {
		_, ok := colourMap[e.Name]
		if !ok {
			if colorIndex >= len(colours) {
				colourMap[e.Name] = "white"
			} else {
				colourMap[e.Name] = colours[colorIndex]
			}
			colorIndex++
		}
	}

	cv := CenterView{
		CenterId:   cw.CenterId,
		CenterName: cw.CenterName,
		Weekdays:   []*WeekdayView{},
	}

	dailyActivities, err := eventsByWeekday(cw.Events)
	if err != nil {
		return nil, err
	}

	spans := []int{}
	gc := 2
	for i := range days {
		sortedActivities := dailyActivities[time.Weekday(i)]

		viewEvents := []*ViewEvent{}
		var prevEvent *ViewEvent
		maxOffset := 0
		for _, a := range sortedActivities {
			e, err := NewViewEvent(a, prevEvent, colourMap)
			if err != nil {
				return nil, err
			}

			viewEvents = append(viewEvents, e)
			if e.Offset > maxOffset {
				maxOffset = e.Offset
			}
			prevEvent = e
		}

		for j, ve := range viewEvents {
			if j+1 < len(viewEvents) {
				nextEvent := viewEvents[j+1]

				if ve.Offset == 0 && ve.prevEvent == nil && nextEvent.Offset == 0 {
					ve.Span = maxOffset + 1
				}
			} else if ve.Offset == 0 && ve.prevEvent == nil {
				ve.Span = maxOffset + 1
			}
		}

		span := maxOffset + 1
		columns := make([]int, span)
		for i := range span {
			columns[i] = gc + i
		}
		cv.Weekdays = append(cv.Weekdays, &WeekdayView{
			Span:        span,
			GridColumns: columns,
			Events:      viewEvents,
		})

		gc += span
		spans = append(spans, span)
	}

	cv.GridColumns = gridColumns(spans)

	return &cv, nil
}

func gridColumns(spans []int) string {
	result := "50px "
	a := lcm(spans)
	for _, s := range spans {
		fr := a / s
		result += strings.Repeat(fmt.Sprintf("%dfr ", fr), s)
	}
	return strings.TrimSpace(result)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func gcd(a, b int) int {
	if min(a, b) == 0 {
		return max(a, b)
	}
	a_1 := max(a, b) % min(a, b)
	return gcd(a_1, min(a, b))
}

func lcm(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	ans := arr[0]

	for i := 1; i < len(arr); i++ {
		ans = (arr[i] * ans) / gcd(arr[i], ans)
	}

	return ans
}

func NewView(plan *CenterPlan) (*View, error) {
	v := View{
		Centers: []*CenterView{},
		Days:    weekdays(),
		Times:   times(),
	}

	for _, p := range plan.Plans {
		cv, err := NewCenterView(p, v.Days)
		if err != nil {
			return nil, err
		}
		v.Centers = append(v.Centers, cv)
	}

	return &v, nil
}
