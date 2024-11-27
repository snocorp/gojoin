package models

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	UNKNOWN = 0
	AM      = 1
	PM      = 2
)

type Activity struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`       // "Swim Creatures 4 - Nigig | Otter"
	Number    string `json:"number"`     // "111222"
	TimeRange string `json:"time_range"` // "9:45 AM - 10:15 AM"
	DetailUrl string `json:"detail_url"`
	DayOfWeek string `json:"days_of_week"` // "Sun"

	startTime *TimeOfDay
	endTime   *TimeOfDay
}

func (a *Activity) Overlaps(other *Activity) (result bool, err error) {
	st, err := a.StartTime()
	if err != nil {
		return
	}

	et, err := a.EndTime()
	if err != nil {
		return
	}

	stOther, err := other.StartTime()
	if err != nil {
		return
	}

	etOther, err := other.EndTime()
	if err != nil {
		return
	}

	overlaps := st.LessThan(&etOther) && stOther.LessThan(&et)

	return overlaps, nil
}

func (a *Activity) EndTime() (TimeOfDay, error) {
	if a.endTime == nil {
		err := a.parseTimeRange()
		if err != nil {
			return TimeOfDay{}, err
		}
	}

	return *a.endTime, nil
}

func (a *Activity) StartTime() (TimeOfDay, error) {
	if a.startTime == nil {
		err := a.parseTimeRange()
		if err != nil {
			return TimeOfDay{}, err
		}
	}

	return *a.startTime, nil
}

func (a *Activity) Weekday() (time.Weekday, error) {
	switch a.DayOfWeek {
	case "Sun":
		return time.Sunday, nil
	case "Mon":
		return time.Monday, nil
	case "Tue":
		return time.Tuesday, nil
	case "Wed":
		return time.Wednesday, nil
	case "Thu":
		return time.Thursday, nil
	case "Fri":
		return time.Friday, nil
	case "Sat":
		return time.Saturday, nil
	}

	return 0, fmt.Errorf("unexpected day of week %v", a.DayOfWeek)
}

func (a *Activity) parseTimeRange() error {
	re := regexp.MustCompile(`^(?:(\d{1,2}):(\d{1,2}) (AM|PM)|(Noon)) - (?:(\d{1,2}):(\d{1,2}) (AM|PM)|(Noon))$`)
	matches := re.FindSubmatch([]byte(a.TimeRange))

	startTimeType := UNKNOWN
	endTimeType := UNKNOWN
	for i, m := range matches {
		if i == 0 {
			continue
		}

		match := string(m)
		if match == "" {
			continue
		}

		if match == "Noon" {
			if i == 4 {
				a.startTime = &TimeOfDay{12, 0}
			} else {
				a.endTime = &TimeOfDay{12, 0}
			}
		} else if match == "AM" {
			if i == 3 {
				startTimeType = AM
			} else {
				endTimeType = AM
			}
		} else if match == "PM" {
			if i == 3 {
				startTimeType = PM
			} else {
				endTimeType = PM
			}
		} else {
			value, err := strconv.ParseInt(match, 10, 8)
			if err != nil {
				return fmt.Errorf("unable to parse match %v in %v", match, a.TimeRange)
			}

			if i == 1 {
				a.startTime = &TimeOfDay{int(value), 0}
			} else if i == 2 {
				a.startTime.Minute = int(value)
			} else if i == 5 {
				a.endTime = &TimeOfDay{int(value), 0}
			} else if i == 6 {
				a.endTime.Minute = int(value)
			} else {
				return fmt.Errorf("unexpected value at %v: %v", i, value)
			}
		}
	}

	if startTimeType == PM && a.startTime.Hour < 12 {
		a.startTime.Hour += 12
	}
	if endTimeType == PM && a.endTime.Hour < 12 {
		a.endTime.Hour += 12
	}

	return nil
}

type ActivitySearchBody struct {
	ActivityItems []*Activity `json:"activity_items"`
}

type ActivityPageInfo struct {
	PageNumber int `json:"page_number"`
	TotalPages int `json:"total_page"`
}

type ActivitySearchHeaders struct {
	PageInfo *ActivityPageInfo `json:"page_info"`
}

type ActivitySearchResponse struct {
	Headers *ActivitySearchHeaders `json:"headers"`
	Body    *ActivitySearchBody    `json:"body"`
}

type ActivitySearchPattern struct {
	Skills                []any    `json:"skills"`                      // [],
	TimeAfter             string   `json:"time_after_str"`              // "",
	DaysOfWeek            string   `json:"days_of_week"`                // "0000000",
	ActivitySelectParam   int      `json:"activity_select_param"`       // 0,
	CenterIds             []string `json:"center_ids"`                  // ["165", "384"],
	TimeBefore            string   `json:"time_before_str"`             // "",
	OpenSpots             *int     `json:"open_spots"`                  // null,
	ActivityId            *string  `json:"activity_id"`                 // null,
	ActivityCategoryIds   []string `json:"activity_category_ids"`       // ["25"],
	DateBefore            string   `json:"date_before"`                 // "",
	MinAge                *int     `json:"min_age"`                     // null,
	DateAfter             string   `json:"date_after"`                  // "",
	ActivityTypeIds       []string `json:"activity_type_ids"`           // ["1"],
	SiteIds               []string `json:"site_ids"`                    // [],
	ForMap                bool     `json:"for_map"`                     // false,
	GeographicAreaIds     []string `json:"geographic_area_ids"`         // [],
	SeasonIds             []string `json:"season_ids"`                  // ["46"],
	ActivityDepartmentIds []string `json:"activity_department_ids"`     // [],
	OtherCategoryIds      []string `json:"activity_other_category_ids"` // [],
	ChildSeasonIds        []string `json:"child_season_ids"`            // [],
	ActivityKeyword       string   `json:"activity_keyword"`            // "otter",
	InstructorIds         []string `json:"instructor_ids"`              // [],
	MaxAge                *int     `json:"max_age"`                     // null,
	CustomPriceFrom       string   `json:"custom_price_from"`           // "",
	CustomPriceTo         string   `json:"custom_price_to"`             // ""
}

type ActivityRequest struct {
	SearchPattern *ActivitySearchPattern `json:"activity_search_pattern"`
}

// func (a *Activity)
