package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	AllowedOrigins       []string      `mapstructure:"ALLOWED_ORIGINS"`
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
	envFiles := []string{
		"app.env",       // Load defaults first
		"app.local.env", // Load local-sensitive overrides
	}

	for _, fileName := range envFiles {
		viper.SetConfigFile(fmt.Sprintf("%s/%s", path, fileName))
		viper.SetConfigType("env")

		if err := viper.MergeInConfig(); err != nil {
			var notFoundErr *viper.ConfigFileNotFoundError
			var pathErr *os.PathError

			if errors.As(err, &notFoundErr) {
				log.Printf("INFO: File %s not found, skipping", fileName)
				continue // File not found, so skip it
			}
			if errors.As(err, &pathErr) && os.IsNotExist(pathErr) {
				log.Printf("INFO: File %s not found, skipping", fileName)
				continue // File not found, so skip it
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

	// Redact sensitive fields
	redactedConfig := config
	redactedConfig.TokenSymmetricKey = "[REDACTED]"
	redactedConfig.EmailSenderPassword = "[REDACTED]"
	redactedConfig.EmailSenderAddress = "[REDACTED]"
	redactedConfig.EmailTestRecipient = "[REDACTED]"
	// redactedConfigJSON, err := json.MarshalIndent(redactedConfig, "", "  ")

	// Use reflection to iterate over struct fields
	v := reflect.ValueOf(redactedConfig)
	t := reflect.TypeOf(redactedConfig)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()
		log.Printf("%s: %v\n", fieldName, fieldValue)
	}

	return config, nil
}
