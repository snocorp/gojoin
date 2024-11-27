package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/snocorp/gojoin/models"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Output the view to HTML",
	Long:  `Render an HTML template using the loaded data.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, err := cmd.Flags().GetString("input")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		inputBytes, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var plan models.Plan
		err = json.Unmarshal(inputBytes, &plan)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		centerMap := map[string]*models.CenterWeek{}
		for _, p := range plan.Plans {
			for _, cw := range p.CenterWeeks {
				centerWeek, ok := centerMap[cw.CenterId]
				if !ok {
					centerWeek = &models.CenterWeek{
						CenterId:   cw.CenterId,
						CenterName: cw.CenterName,
						Events:     []*models.Activity{},
					}
					centerMap[cw.CenterId] = centerWeek
				}

				centerWeek.Events = append(centerWeek.Events, cw.Events...)
			}
		}

		centerPlan := &models.CenterPlan{Plans: []*models.CenterWeek{}}
		for _, cw := range centerMap {
			centerPlan.Plans = append(centerPlan.Plans, cw)
		}

		view, err := models.NewView(centerPlan)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		funcMap := template.FuncMap{
			"css": func(s string) template.CSS {
				return template.CSS(s)
			},
		}

		filename := "./templates/week.html.gotmpl"
		name := path.Base(filename)
		tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = tmpl.Execute(os.Stdout, view)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly:
	viewCmd.Flags().String("input", "", "The input file to load into the view")
}
