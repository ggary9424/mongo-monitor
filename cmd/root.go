package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mongoURI string = "mongodb://127.0.0.1:27017"

var rootCmd = &cobra.Command{
	Use:   "mongo-monitor",
	Short: "Mongo monitor",
	Long:  "It is a mongo monitor",
}

// Execute represents executing cobra library
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()

	pf.Bool("debug", false, "Run the program with debug mode")
	pf.StringVar(&mongoURI, "uri", "mongodb://127.0.0.1:27017", "URI of mongo you want to monitor")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	cobra.MarkFlagRequired(pf, "uri")
}
