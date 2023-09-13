package config

import "github.com/spf13/viper"

type Config struct {
	DB_Host         string `mapstructure:"DB_HOST"`
	DB_Password     string `mapstructure:"DB_PASSWORD"`
	DB_User         string `mapstructure:"DB_USER"`
	DB_Name         string `mapstructure:"DB_NAME"`
	DB_SSLMode      string `mapstructure:"DB_SSL_MODE"`
	DB_Port         string `mapstructure:"DB_PORT"`
	Token_Secret    string `mapstructure:"TOKEN_SECRET"`
	Email_From      string `mapstructure:"EMAIL_FROM"`
	Email_Smtp_Host string `mapstructure:"EMAIL_SMTP_HOST"`
	Email_Smtp_User string `mapstructure:"EMAIL_SMTP_USER"`
	Email_Smtp_Pass string `mapstructure:"EMAIL_SMTP_PASS"`
	Email_Smtp_Port int    `mapstructure:"EMAIL_SMTP_PORT"`
	Client_Origin   string `mapstructure:"CLIENT_ORIGIN"`
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
