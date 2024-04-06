package cmd

import (
	"fmt"
	"github.com/calamity-m/paternity/pkg/config"
	"github.com/calamity-m/paternity/pkg/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "paternity",
	Short:   "App for testing if a container image is built from a parent image",
	Long:    ``,
	Version: version.GetVersionS(),
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
	// On startup of Execute run initConfig
	cobra.OnInitialize(initialize)

	// Flags
	rootCmd.PersistentFlags().StringVar(&cfg,
		"config",
		"",
		"config file (default is $HOME/.config/paternity.yaml)")
}

func initialize() {
	// Setup viper
	initConfig()

	// Set zerolog
	initLog()

	// Log some default information
	log.Debug().Msgf("Config file used: %s", viper.ConfigFileUsed())
}

func initConfig() {

	viper.SetDefault("log.level", zerolog.LevelInfoValue)
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

		// Search config in home directory with name ".paternity" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("paternity")
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
			_, _ = fmt.Fprintln(os.Stderr, "Failed to write initial config file to", cfgDir, "/paternity.yaml")
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

	// Initialize log level
	level, err := zerolog.ParseLevel(configuration.Log.Level)
	if err != nil {
		return
	}
	zerolog.SetGlobalLevel(level)

	// Initialize log format
	if !configuration.Log.Structured {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Debug().Msg("Initialized logging")
	log.Debug().Msgf("Logging level: %s", configuration.Log.Level)
	log.Debug().Msgf("Logging structured: %t", configuration.Log.Structured)
}
