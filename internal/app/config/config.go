package config

import (
	v0 "supermarket/internal/adapter/http/v0"
	"supermarket/pkg/database/postgres"

	"supermarket/pkg/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Logger   logrus.Config
	Postgres postgres.Config
	Web      v0.Config
}

func GetConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := bindEnvConfig(); err != nil {
		panic(errors.Wrap(err, "error reading env"))
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "error reading config file"))
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(errors.Wrap(err, "error unmarshaling config file"))
	}
	return cfg
}

func bindEnvConfig() error {
	bindings := make(map[string]string)
	bindings["web.host"] = "HOST"
	bindings["web.port"] = "PORT"

	bindings["postgres.database"] = "POSTGRES_DATABASE"
	bindings["postgres.host"] = "POSTGRES_HOST"
	bindings["postgres.port"] = "POSTGRES_PORT"
	bindings["postgres.user"] = "POSTGRES_USER"
	bindings["postgres.password"] = "POSTGRES_PASSWORD"

	for name, binding := range bindings {
		if err := viper.BindEnv(name, binding); err != nil {
			return err
		}
	}

	return nil
}
