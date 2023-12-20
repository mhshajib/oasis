package cmd

import (
	"fmt"
	"oasis/pkg/config"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// rootCmd is the root command of backup service
	rootCmd = &cobra.Command{
		Use:   "oasis",
		Short: "oasis service provide cli for making clean block function easier",
		Long:  `oasis service provide cli for making clean block function easier`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Loading configurations")
	if err := config.Init(); err != nil {
		logrus.Warn("Failed to load configuration")
		logrus.Fatal(err)
	}
	logrus.Info("Configurations loaded successfully!")
}
