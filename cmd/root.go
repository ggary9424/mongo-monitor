package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgPath string

var rootCmd = &cobra.Command{
	Use:   "mongo-monitor",
	Short: "Mongo monitor",
	Long:  "It is a mongo monitor",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file path")
	cobra.MarkFlagRequired(pf, "config")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("Must to assign a config file path")
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// Initialize logger system
	debug := viper.GetBool("system.debug")

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stdout)
}

func setDefaultConfig() {
	// system
	viper.SetDefault("system.debug", "false")

	// mongo
	viper.SetDefault("mongo.uri", "mongodb://127.0.0.1:27017")
}
