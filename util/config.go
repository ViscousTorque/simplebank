package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`      // Matches DB_DRIVER in app.env
	DBSource      string `mapstructure:"DB_SOURCE"`      // Matches DB_SOURCE in app.env
	ServerAddress string `mapstructure:"SERVER_ADDRESS"` // Matches SERVER_ADDRESS in app.env
}

// loadConfig - read from file or env vars
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
