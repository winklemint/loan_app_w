package con

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`

		Name string `yaml:"name"`
	} `yaml:"database"`

	HTTP struct {
		Port int `yaml:"port"`
	} `yaml:"http"`
}

var db *sql.DB

func ConnectDB(dbConfig Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:@tcp(%s:%d)/%s", dbConfig.Database.Username, dbConfig.Database.Host, dbConfig.Database.Port, dbConfig.Database.Name)
	var err error
	db, err = sql.Open(dbConfig.Database.Driver, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func LoadConfig(file string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func GetDB() (*sql.DB, error) {
	config, err := LoadConfig("Config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	db, err := ConnectDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
