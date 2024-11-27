package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/snocorp/gojoin/models"
)

type GetFiltersOptions struct {
	NoCache bool
	Verbose bool
}

func GetFilters(options GetFiltersOptions) (body models.FiltersBody, err error) {
	loadedCachedFilter := false
	useCache := !options.NoCache

	var filterBytes []byte
	var filters models.FiltersResponse
	if useCache {
		wd, err := os.Getwd()
		if err != nil {
			return body, err
		}

		cacheDir := path.Join(wd, ".gojoin", "cache")

		filterBytes, err = os.ReadFile(path.Join(cacheDir, "filters.json"))
		if err != nil {
			if options.Verbose {
				fmt.Println(err)
			}
		} else {
			err = json.Unmarshal(filterBytes, &filters)
			if err != nil {
				if options.Verbose {
					fmt.Println(err)
				}
			} else {
				loadedCachedFilter = true
			}
		}
	}

	if !loadedCachedFilter {
		now := time.Now().UnixMilli()
		resp, err := http.Get(fmt.Sprintf("https://anc.ca.apm.activecommunities.com/ottawa/rest/activities/filters?locale=en-US&ui_random=%v", now))
		if err != nil {
			return models.FiltersBody{}, err
		}
		defer resp.Body.Close()
		filterBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return models.FiltersBody{}, err
		}

		err = json.Unmarshal(filterBytes, &filters)
		if err != nil {
			return models.FiltersBody{}, err
		}

		if useCache {
			wd, err := os.Getwd()
			if err != nil {
				return body, err
			}

			cacheDir := path.Join(wd, ".gojoin", "cache")

			err = os.MkdirAll(cacheDir, 0775)
			if err != nil {
				if options.Verbose {
					fmt.Println(err)
				}
			} else {
				err = os.WriteFile(path.Join(cacheDir, "filters.json"), filterBytes, 0664)
				if err != nil && options.Verbose {
					fmt.Println(err)
				}
			}
		}
	}

	return filters.Body, nil
}
