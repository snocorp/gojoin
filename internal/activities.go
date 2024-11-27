package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/snocorp/gojoin/models"
)

type GetActivitiesOptions struct {
	Verbose bool
}

// {"order_by":"Name","page_number":2,"total_records_per_page":20}
type PageInfo struct {
	OrderBy    string `json:"order_by"`
	PageNumber int    `json:"page_number"`
	PerPage    int    `json:"total_records_per_page"`
}

func GetActivities(request models.ActivityRequest, options GetActivitiesOptions) (activities []*models.Activity, err error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return
	}

	morePages := true
	pageNumber := 1
	var result []*models.Activity
	for morePages {
		result, morePages, err = getActivities(requestBytes, pageNumber)
		if err != nil {
			return activities, err
		}

		activities = append(activities, result...)
		pageNumber += 1
	}

	return
}

func getActivities(requestBytes []byte, page int) ([]*models.Activity, bool, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(
		"POST",
		"https://anc.ca.apm.activecommunities.com/ottawa/rest/activities/list?locale=en-US",
		bytes.NewReader(requestBytes),
	)
	if err != nil {
		return []*models.Activity{}, false, err
	}

	pageInfo := PageInfo{
		OrderBy:    "Name",
		PerPage:    20,
		PageNumber: page,
	}
	pageInfoJson, err := json.Marshal(pageInfo)
	if err != nil {
		return []*models.Activity{}, false, err
	}

	req.Header.Add("Page_info", string(pageInfoJson))
	req.Header.Add("content-length", strconv.FormatInt(int64(len(requestBytes)), 10))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	req.Header.Add("content-type", "application/json;charset=utf-8")
	req.Header.Add("origin", "https://anc.ca.apm.activecommunities.com")
	req.Header.Add("x-csrf-token", uuid.NewString())

	resp, err := client.Do(req)
	if err != nil {
		return []*models.Activity{}, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []*models.Activity{}, false, fmt.Errorf("unexpected status %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []*models.Activity{}, false, err
	}

	var searchResponse models.ActivitySearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		fmt.Println(string(body))
		return []*models.Activity{}, false, err
	}

	morePages := searchResponse.Headers.PageInfo.PageNumber < searchResponse.Headers.PageInfo.TotalPages
	return searchResponse.Body.ActivityItems, morePages, nil
}
