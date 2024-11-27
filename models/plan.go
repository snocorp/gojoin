package models

type CenterWeek struct {
	CenterId   string      `json:"center_id"`
	CenterName string      `json:"center_name"`
	Events     []*Activity `json:"events"`
}

type PersonCenterWeek struct {
	Person      string        `json:"person"`
	CenterWeeks []*CenterWeek `json:"center_weeks"`
}

type Plan struct {
	Plans []*PersonCenterWeek `json:"plans"`
}

type CenterPlan struct {
	Plans []*CenterWeek `json:"plans"`
}
