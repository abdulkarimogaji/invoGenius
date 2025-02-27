package config

import (
	"github.com/spf13/viper"
)

type config struct {
	Port         int    `mapstructure:"PORT"`
	Database_Uri string `mapstructure:"DATABASE_URI"`
	Token_Secret string `mapstructure:"TOKEN_SECRET"`
	App_Uri      string `mapstructure:"APP_URI"`
	Token_Expire int    `mapstructure:"TOKEN_EXPIRE"`
}

var C config

func LoadConfig() error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&C)
	return err
}
