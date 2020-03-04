package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProcessEnv  string `envconfig:"ISSUINGSERVICE_PROCESS_ENV"`
	Addr        string `envconfig:"ISSUINGSERVICE_SERVER_ADDR"`
	LogFilename string `envconfig:"ISSUINGSERVICE_LOGS_FILENAME"`
}

// Env returns the settings from the environment
func Env() (conf Config, err error) {
	err = envconfig.Process("", &conf)
	if err != nil {
		return
	}

	if conf.ProcessEnv == "" {
		conf.ProcessEnv = "dev"
	}

	if conf.Addr == "" {
		conf.Addr = ":8082"
	}

	if conf.LogFilename == "" {
		conf.LogFilename = "issuing_service.log"
	}

	return
}
