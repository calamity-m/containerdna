package cmd

import (
	"fmt"
	"github.com/calamity-m/containerdna/pkg/config"
	"github.com/calamity-m/containerdna/pkg/version"
	"github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "containerdna",
	Short:   "Simple CLI for testing container heritage",
	Long:    `Compares container layer blocks to assert ancestry of containers`,
	Version: version.GetVersionS(),
}

var heritageGroup = &cobra.Group{
	ID:    "heritage",
	Title: "Heritage",
}

var configuration *config.Config

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Set up groups
	rootCmd.AddGroup(heritageGroup)

	// On startup of Execute run initConfig
	cobra.OnInitialize(initialize)

	// Flags
	rootCmd.PersistentFlags().StringVar(&cfg,
		"config",
		"",
		"config file (default is $HOME/.config/containerdna.yaml)")
}

func initialize() {
	// Setup viper
	initConfig()

	// Set logging
	initLog()

	// Log some default information
	logrus.Debugf("Config file used: %s", viper.ConfigFileUsed())
}

func initConfig() {

	viper.SetDefault("log.level", logrus.InfoLevel)
	viper.SetDefault("log.file", "")
	viper.SetDefault("log.structured", false)

	var cfgDir string
	// Define configuration file
	if cfg != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg)
	} else {
		// Find home directory.
		var err error
		cfgDir, err = os.UserConfigDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".containerdna" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("containerdna")
	}

	// Prepend PATERNITY to all env vars
	viper.SetEnvPrefix("PATERNITY")
	// Replace all hyphens with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		err = viper.SafeWriteConfig()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to write initial config file to", cfgDir, "/containerdna.yaml")
			return
		}
	}
}

func initLog() {
	var err error
	configuration, err = config.New()

	if err != nil {
		return
	}

	// Initalize logrus log level
	lvl, err := logrus.ParseLevel(configuration.Log.Level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Failed to parse log level, defaulting to debug")
	} else {
		logrus.SetLevel(lvl)
	}

	// Initialize log structure
	if configuration.Log.Structured {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// Initialize log file
	if configuration.Log.File != "" {
		var file, err = os.OpenFile(configuration.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.Error("Failed to open log file, defaulting to stderr")
		} else {
			logrus.SetOutput(file)
		}
	}
}
