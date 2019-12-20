package configs

import (
	"log"

	"github.com/caarlos0/env"
)

func ReadConfig(cfgPtr interface{}) error {
	return env.Parse(cfgPtr)
}

func MustReadConfig(cfgPtr interface{}) {
	err := ReadConfig(cfgPtr)
	if err != nil {
		log.Fatalf("parse configs: %v", err)
	}
}
