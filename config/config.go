package config

import (
	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	viper.BindEnv("PORT")
	viper.BindEnv("DATABASE_URI")
	viper.BindEnv("TOKEN_SECRET")
	viper.BindEnv("APP_URI")
	viper.BindEnv("TOKEN_EXPIRE")

	viper.AutomaticEnv()
	err = viper.Unmarshal(&C)
	return err
}
