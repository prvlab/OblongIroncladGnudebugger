package config

import (
	"encoding/json"
	"os"
)

// Config структура конфигурации
type Config struct {
	Telegram struct {
		Token string `json:"token"`
	} `json:"telegram"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DbName   string `json:"dbname"`
	} `json:"database"`
	Server struct {
		ServerHost string `json:"server_host"`
		ServerPort string `json:"server_port"`
	}
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}