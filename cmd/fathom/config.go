package main

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/usefathom/fathom/pkg/datastore"
)

type Config struct {
	Database *datastore.Config

	Secret string
}

func parseConfig() *Config {
	var cfg Config
	godotenv.Load()
	err := envconfig.Process("Fathom", &cfg)
	if err != nil {
		log.Fatalf("Error parsing Fathom config from environment: %s", err)
	}

	// alias sqlite to sqlite3
	if cfg.Database.Driver == "sqlite" {
		cfg.Database.Driver = "sqlite3"
	}

	// if secret key is empty, use a randomly generated one to ease first-time installation
	if cfg.Secret == "" {
		cfg.Secret = randomString(40)
		os.Setenv("FATHOM_SECRET", cfg.Secret)
	}

	return &cfg
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}

	return string(bytes)
}
