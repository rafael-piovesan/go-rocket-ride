package rocketride

import (
	"fmt"
	"log"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	cfgFileName = "app"
	cfgFileType = "env"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	IdemKeyTimeout  int    `mapstructure:"IDEM_KEY_TIMEOUT" validate:"required"`
	DBSource        string `mapstructure:"DB_SOURCE"  validate:"required"`
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"  validate:"required"`
	StripeKey       string `mapstructure:"STRIPE_KEY"  validate:"required"`
	DatastoreAccess string `mapstructure:"DATASTORE_ACCESS" validate:"required,oneof=bun sqlc jet"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(cfgPath string) (config Config, err error) {
	// set config file path, name and extension (e.g,'/path-to-config/app.env')
	viper.AddConfigPath(cfgPath)
	viper.SetConfigName(cfgFileName)
	viper.SetConfigType(cfgFileType)

	// bind env vars, to be used in case the config file is not found
	_ = viper.BindEnv("IDEM_KEY_TIMEOUT")
	_ = viper.BindEnv("DB_SOURCE")
	_ = viper.BindEnv("SERVER_ADDRESS")
	_ = viper.BindEnv("STRIPE_KEY")

	// default config values
	viper.SetDefault("IDEM_KEY_TIMEOUT", 5)
	viper.SetDefault("DATASTORE_ACCESS", "bun")

	// enable loading of env vars values
	// it'll override the values in the config file
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("config file '%s' not found", path.Join(cfgPath, fmt.Sprintf("%s.%s", cfgFileName, cfgFileType)))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	v := validator.New()
	err = v.Struct(&config)
	return config, err
}
