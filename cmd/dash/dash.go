package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ricoberger/dash/pkg/dashboard"
	"github.com/ricoberger/dash/pkg/datasource"
	fLog "github.com/ricoberger/dash/pkg/log"
	"github.com/ricoberger/dash/pkg/render"
	"github.com/ricoberger/dash/pkg/version"

	"github.com/spf13/cobra"
)

var (
	configDir      string
	configInterval string
	configRefresh  string
	debug          bool
)

var rootCmd = &cobra.Command{
	Use:   "dash",
	Short: "dash - terminal dashboard.",
	Long:  "dash - terminal dashboard.",
	Run: func(cmd *cobra.Command, args []string) {
		if configDir == "~/.dash" {
			configDir = os.Getenv("HOME") + "/.dash"
		}

		err := fLog.Init(configDir, debug)
		if err != nil {
			log.Fatalf("Could not open log file: %v", err)
		}
		defer fLog.Close()

		datasources, err := datasource.New(configDir)
		if err != nil {
			log.Fatalf("Could not load datasources: %v", err)
		}

		dashboards, err := dashboard.New(configDir)
		if err != nil {
			log.Fatalf("Could not load dashboards: %v", err)
		}

		err = render.Run(datasources, dashboards, configInterval, configRefresh)
		if err != nil {
			log.Fatalf("Unexpected error: %v", err)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information for dash.",
	Long:  "Print version information for dash.",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := version.Print("dash")
		if err != nil {
			log.Fatalf("Failed to print version information: %v", err)
		}

		fmt.Fprintln(os.Stdout, v)
		return
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configDir, "config.dir", "~/.dash", "Location of the datasources and dashboards folder.")
	rootCmd.PersistentFlags().StringVar(&configInterval, "config.interval", "1h", "Interval to retrieve data for.")
	rootCmd.PersistentFlags().StringVar(&configRefresh, "config.refresh", "5m", "Time between refreshs of the dashboard.")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Log debug information.")

	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to initialize dash: %v", err)
	}
}
