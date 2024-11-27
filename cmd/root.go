package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var verbose bool
var noCache bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gojoin",
	Short: "Tool for visualizing weekly activities",
	Long:  `A tool to load and then view activities in a weekly format.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Output verbose information")
	rootCmd.PersistentFlags().BoolVar(&noCache, "nocache", false, "Do not use cached data")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
