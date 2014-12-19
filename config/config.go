package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v1"
)

var DatabaseYAMLPath = os.Getenv("GOPATH") + "/src/github.com/fabioberger/dockerize-tutorial/db/dbconf.yml"

var Env string = os.Getenv("GO_ENV")

var Database dbConfig

type config struct {
	Database dbConfig
}

type dbConfig struct {
	Driver string
	Open   string
}

func Init() {
	if Env == "production" {
		ParseDatabaseYAML("production")
	} else if Env == "development" || Env == "" {
		ParseDatabaseYAML("development")
	} else if Env == "test" {
		ParseDatabaseYAML("test")
	} else {
		panic("Unkown environment. Don't know what configuration to use!")
	}
}

func ParseDatabaseYAML(env string) {
	// Read all data from the file and unmarshall it
	var data map[string]struct {
		Driver string
		Open   string
	}
	content, err := ioutil.ReadFile(DatabaseYAMLPath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(content, &data); err != nil {
		panic(err)
	}
	envData := data[env]

	// parse the env variable and set it properly
	if strings.Contains(envData.Open, "$POSTGRES_PORT_5432_TCP_ADDR") {
		envData.Open = strings.Replace(envData.Open, "$POSTGRES_PORT_5432_TCP_ADDR", os.Getenv("POSTGRES_PORT_5432_TCP_ADDR"), -1)
	}
	fmt.Println("[database] Using psql paramaters:", envData.Open)

	// Construct a dbConfig object from envData
	Database = dbConfig{
		Driver: envData.Driver,
		Open:   envData.Open,
	}
}
