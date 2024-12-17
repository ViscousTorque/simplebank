package util

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"` // Matches DB_DRIVER in app.env
	DBSource             string        `mapstructure:"DB_SOURCE"` // Matches DB_SOURCE in app.env
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	HttpServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"` // Matches HTTP_SERVER_ADDRESS in app.env
	GrpcServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"` // Matches HTTP_SERVER_ADDRESS in app.env
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EnableReflection     bool          `mapstructure:"ENABLE_REFLECTION"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	EmailTestRecipient   string        `mapstructure:"EMAIL_TEST_RECIPIENT"`
}

// LocalConfig - setting app.env.local file for local testing of email sending using sensitive data
func LoadConfig(path string) (config Config, err error) {
	// List of files to load, in the order of precedence
	envFiles := []string{
		"app.env",       // Load defaults first
		"app.local.env", // Load local-sensitive overrides (from GitHub CI/CD)
	}

	for _, fileName := range envFiles {
		viper.SetConfigFile(fmt.Sprintf("%s/%s", path, fileName))
		viper.SetConfigType("env") // Explicitly set the type as "env" for .env files

		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Printf("File %s not found, skipping", fileName)
				continue
			} else {
				return config, fmt.Errorf("error loading %s: %w", fileName, err)
			}
		} else {
			log.Printf("Loaded config file: %s", fileName)
		}
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return config, nil
}
