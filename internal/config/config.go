package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	AppConfig
	DBConfig
}

type AppConfig struct {
	Port string
}

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
	Sslmode  string
}

func Load() (*Config, error) {
	var appCfg AppConfig
	var dbCfg DBConfig
	var err error

	if appCfg.Port, err = getStrValue("APP_PORT", true); err != nil {
		return nil, err
	}

	if dbCfg.User, err = getStrValue("DB_USER", true); err != nil {
		return nil, err
	}
	if dbCfg.Password, err = getStrValue("DB_PASSWORD", true); err != nil {
		return nil, err
	}
	if dbCfg.Host, err = getStrValue("DB_HOST", true); err != nil {
		return nil, err
	}
	if dbCfg.Port, err = getIntValue("DB_PORT", true); err != nil {
		return nil, err
	}
	if dbCfg.Name, err = getStrValue("DB_NAME", true); err != nil {
		return nil, err
	}
	if dbCfg.Sslmode, err = getStrValue("DB_SSLMODE", true); err != nil {
		return nil, err
	}

	return &Config{
		AppConfig: appCfg,
		DBConfig:  dbCfg,
	}, nil
}

func getStrValue(key string, required bool) (string, error) {
	value := os.Getenv(key)
	if key == "" && required {
		return "", errors.New("key is empty")
	}
	return value, nil
}

func getIntValue(key string, required bool) (int, error) {
	strValue, err := getStrValue(key, required)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}
	return value, nil
}
