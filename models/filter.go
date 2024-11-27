package models

type Criterium struct {
	Id          string `json:"id"`
	Description string `json:"desc"`
}

type Criteria = []Criterium

type FiltersBody struct {
	Centers    Criteria `json:"centers"`
	Categories Criteria `json:"categories"`
	Seasons    Criteria `json:"seasons"`
}

type FiltersResponse struct {
	Body FiltersBody `json:"body"`
}
