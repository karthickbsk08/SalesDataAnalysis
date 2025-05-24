package dbconnection

import (
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type AllDatabases struct {
	MariaDB     *DBConfig
	PostGresSQL *DBConfig
}

type DBConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	DBType   string
	DBflag   string
}

func (db *AllDatabases) Init() {
	db.PostGresSQL = LoadDBDetailsfromEnv("dev")
	db.PostGresSQL.DBflag = "POSTGRESSQLDB"
	// db.MariaDB=LoadDBDetailsfromEnv("test")
}

// LoadEnvToMap loads the .env file into a map without setting OS environment variables
func LoadEnvToMap(envFile string) (map[string]string, error) {
	envMap, err := godotenv.Read(envFile)
	if err != nil {
		return nil, err
	}
	return envMap, nil
}

func LoadDBConfigFromMap(envMap map[string]string) (*DBConfig, error) {
	portStr, ok := envMap["DB_PORT"]
	if !ok {
		return nil, fmt.Errorf("DB_PORT not found")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return &DBConfig{
		Server:   envMap["DB_SERVER"],
		Port:     port,
		User:     envMap["DB_USER"],
		Password: envMap["DB_PASSWORD"],
		Database: envMap["DB_DATABASE"],
		DBType:   envMap["DB_TYPE"],
	}, nil
}

func LoadDBDetailsfromEnv(env string) *DBConfig {

	envFile := ".env." + env

	envMap, err := LoadEnvToMap(envFile)
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	config, err := LoadDBConfigFromMap(envMap)
	if err != nil {
		log.Fatalf("Failed to load DB config: %v", err)
	}

	return config
}
