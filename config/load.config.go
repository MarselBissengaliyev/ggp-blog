package config

import "github.com/spf13/viper"

type Config struct {
	DB_Host      string `mapstructure:"DB_HOST"`
	DB_Password  string `mapstructure:"DB_PASSWORD"`
	DB_User      string `mapstructure:"DB_USER"`
	DB_Name      string `mapstructure:"DB_NAME"`
	DB_SSLMode   string `mapstructure:"DB_SSL_MODE"`
	DB_Port      string `mapstructure:"DB_PORT"`
	Token_Secret string `mapstructure:"TOKEN_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}
