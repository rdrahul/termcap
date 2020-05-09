package cmd

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	utils "github.com/rdrahul/termcap/utils"
	"github.com/spf13/viper"
)

//basic configuration setting
var (
	Application = "termcap"
	Version     = "0.0.1"
)

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			utils.Er(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".termcap")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
