package config

import (
	"math/rand"
	"os"
	"path/filepath"
	"net/url"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/usefathom/fathom/pkg/datastore/sqlstore"
)

// Config wraps the configuration structs for the various application parts
type Config struct {
	Database *sqlstore.Config
	Secret   string
}

// LoadEnv loads env values from the supplied file
func LoadEnv(file string) {
	if file == "" {
		log.Warn("Missing configuration file. Using defaults.")
		return
	}

	absFile, _ := filepath.Abs(file)
	_, err := os.Stat(absFile)
	fileNotExists := os.IsNotExist(err)

	if fileNotExists {
		log.Warnf("Error reading configuration. File `%s` does not exist.", file)
		return
	}

	log.Printf("Configuration file: %s", absFile)

	// read file into env values
	err = godotenv.Load(absFile)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %s", err)
	}
}

// Parse environment into a Config struct
func Parse() *Config {
	var cfg Config

	// with config file loaded into env values, we can now parse env into our config struct
	err := envconfig.Process("Fathom", &cfg)
	if err != nil {
		log.Fatalf("Error parsing configuration from environment: %s", err)
	}

	if cfg.Database.URL != "" {
		u, err := url.Parse(cfg.Database.URL)
		if err != nil {
			log.Fatalf("Error parsing DATABASE_URL from environment: %s", err)
		}
		cfg.Database.Driver = u.Scheme
	}

	// alias sqlite to sqlite3
	if cfg.Database.Driver == "sqlite" {
		cfg.Database.Driver = "sqlite3"
	}

	// use absolute path to sqlite3 database
	if cfg.Database.Driver == "sqlite3" {
		cfg.Database.Name, _ = filepath.Abs(cfg.Database.Name)
	}

	// if secret key is empty, use a randomly generated one
	if cfg.Secret == "" {
		cfg.Secret = randomString(40)
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
