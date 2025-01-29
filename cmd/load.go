package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/snocorp/gojoin/internal"
	"github.com/snocorp/gojoin/models"
	"github.com/spf13/cobra"
)

type LoadOptions struct {
	season       models.Criterium
	center       models.Criterium
	category     models.Criterium
	searchString string
	outputPath   string
	person       string
	verbose      bool
}

func getOptions(cmd *cobra.Command) (*LoadOptions, error) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return nil, err
	}

	options := internal.GetFiltersOptions{
		Verbose: verbose,
	}
	filters, err := internal.GetFilters(options)
	if err != nil {
		return nil, err
	}

	seasonId, err := cmd.Flags().GetString("season")
	if err != nil {
		return nil, err
	}

	centerId, err := cmd.Flags().GetString("center")
	if err != nil {
		return nil, err
	}

	categoryId, err := cmd.Flags().GetString("category")
	if err != nil {
		return nil, err
	}

	searchString, err := cmd.Flags().GetString("search")
	if err != nil {
		return nil, err
	}

	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return nil, err
	}

	person, err := cmd.Flags().GetString("person")
	if err != nil {
		return nil, err
	}

	season, err := promptSeason(filters, seasonId)
	if err != nil {
		return nil, err
	}

	center, err := promptCenter(filters, centerId)
	if err != nil {
		return nil, err
	}

	category, err := promptCategory(filters, categoryId)
	if err != nil {
		return nil, err
	}

	return &LoadOptions{
		season:       season,
		center:       center,
		category:     category,
		searchString: searchString,
		outputPath:   outputPath,
		person:       person,
		verbose:      verbose,
	}, nil
}

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load data for actitvities",
	Long:  `Requests data using the given search criteria and stores it in the output file.`,
	Run: func(cmd *cobra.Command, args []string) {
		options, err := getOptions(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		req := models.ActivityRequest{
			SearchPattern: &models.ActivitySearchPattern{
				SeasonIds:           []string{options.season.Id},
				CenterIds:           []string{options.center.Id},
				ActivityCategoryIds: []string{options.category.Id},
				ActivityKeyword:     options.searchString,
			},
		}

		activities, err := internal.GetActivities(req, internal.GetActivitiesOptions{
			Verbose: options.verbose,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if options.outputPath == "" {
			options.outputPath = "output.json"
		}

		var existingPlan models.Plan
		plan := models.Plan{Plans: []*models.PersonCenterWeek{}}
		outputBytes, err := os.ReadFile(options.outputPath)
		if err != nil {
			if options.verbose {
				fmt.Println("unable to read output file to load existing data")
			}
		}

		err = json.Unmarshal(outputBytes, &existingPlan)
		if err != nil {
			if options.verbose {
				fmt.Println("unable to unmarshal existing data")
			}
		}

		if options.verbose {
			fmt.Printf("Found %v activities\n", len(activities))
		}

		var personWeek *models.PersonCenterWeek
		foundPlan := false
		for _, p := range existingPlan.Plans {

			if p.Person == options.person {
				if options.verbose {
					fmt.Printf("Found plan for %v\n", options.person)
				}

				foundCenterWeek := false
				for _, cw := range p.CenterWeeks {
					if cw.CenterId == options.center.Id {
						cw.Events = activities
						foundCenterWeek = true
					}
				}
				if !foundCenterWeek {
					p.CenterWeeks = append(p.CenterWeeks, &models.CenterWeek{
						CenterId:   options.center.Id,
						CenterName: options.center.Description,
						Events:     activities,
					})
				}
				foundPlan = true
			}

			plan.Plans = append(plan.Plans, p)
		}
		if !foundPlan {
			if options.verbose {
				fmt.Printf("Creating plan for %v\n", options.person)
			}
			personWeek = &models.PersonCenterWeek{Person: options.person, CenterWeeks: []*models.CenterWeek{
				{CenterId: options.center.Id, CenterName: options.center.Description, Events: activities},
			}}
			plan.Plans = append(plan.Plans, personWeek)
		}

		if options.verbose {
			fmt.Printf("Found %v activities\n", len(plan.Plans[0].CenterWeeks[0].Events))
		}

		planJson, err := json.Marshal(plan)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = os.WriteFile(options.outputPath, planJson, 0664)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	loadCmd.Flags().String("season", "", "The season ID")
	loadCmd.Flags().String("center", "", "The center ID")
	loadCmd.Flags().String("category", "", "The category ID")
	loadCmd.Flags().String("search", "", "The search string")
	loadCmd.Flags().String("output", "", "The output file for the loaded data")

	loadCmd.Flags().String("person", "", "The person with which the events will be associated")
	loadCmd.MarkFlagRequired("person")
}

func promptSeason(filters models.FiltersBody, defaultId string) (models.Criterium, error) {
	if defaultId != "" {
		for _, s := range filters.Seasons {
			if s.Id == defaultId {
				return s, nil
			}
		}
	}

	seasonItems := []string{}
	for _, s := range filters.Seasons {
		seasonItems = append(seasonItems, fmt.Sprintf("%v (%v)", s.Description, s.Id))
	}

	prompt := promptui.Select{
		Label: "Seasons",
		Items: seasonItems,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return models.Criterium{}, err
	}

	return filters.Seasons[index], nil
}

func promptCenter(filters models.FiltersBody, defaultId string) (models.Criterium, error) {
	if defaultId != "" {
		for _, c := range filters.Centers {
			if c.Id == defaultId {
				return c, nil
			}
		}
	}

	centerItems := []string{}
	for _, c := range filters.Centers {
		centerItems = append(centerItems, fmt.Sprintf("%v (%v)", c.Description, c.Id))
	}

	prompt := promptui.Select{
		Label: "Center",
		Items: centerItems,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return models.Criterium{}, err
	}

	return filters.Centers[index], nil
}

func promptCategory(filters models.FiltersBody, defaultId string) (models.Criterium, error) {
	if defaultId != "" {
		for _, c := range filters.Categories {
			if c.Id == defaultId {
				return c, nil
			}
		}
	}

	categoryItems := []string{}
	for _, c := range filters.Categories {
		categoryItems = append(categoryItems, fmt.Sprintf("%v (%v)", c.Description, c.Id))
	}

	prompt := promptui.Select{
		Label: "Category",
		Items: categoryItems,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return models.Criterium{}, err
	}

	return filters.Categories[index], nil
}
