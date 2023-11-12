package config

import "github.com/spf13/viper"

type Config struct {
	Port  string `mapstructure:"PORT"`
	DBUrl string `mapstructure:"DB_URL"`
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"`
}

func LoadConfig() (c *Config, err error) {

	viper.AddConfigPath("./configs/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	// if the default files could not be set, we define some default values
	if err != nil {
		viper.SetDefault("PORT", "3000")
		viper.SetDefault("DB_URL", "./last_resort.db")
	}

	err = viper.Unmarshal(&c)

	return
}
