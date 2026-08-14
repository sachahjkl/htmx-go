package config

import "github.com/spf13/viper"

type Config struct {
	Port          string `mapstructure:"PORT"`
	DBUrl         string `mapstructure:"DB_URL"`
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"`
	Version       string `mapstructure:"VERSION"`
	CommitSha     string `mapstructure:"COMMIT_SHA"`
}

func LoadConfig() (c *Config, err error) {

	viper.AutomaticEnv()

	// need to manually bind otherwise values are empty (why ?)
	viper.BindEnv("VERSION")
	viper.BindEnv("COMMIT_SHA")

	err = viper.ReadInConfig()

	// if the default files could not be set, we define some default values
	if err != nil {
		viper.SetDefault("PORT", "3000")
		viper.SetDefault("DB_URL", "./last_resort.db")
		viper.SetDefault("ENCRYPTION_KEY", "ABCDEFGHJKLMNOPQRSTUVWXYZ")
		viper.SetDefault("VERSION", "v0.0.0+empty")
		viper.SetDefault("COMMIT_SHA", "abcdefg")
	}

	err = viper.Unmarshal(&c)

	return
}
