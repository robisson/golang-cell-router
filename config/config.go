package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Cells []CellConfig `json:"cells"`
}

type CellConfig struct {
	Name      string `json:"name"`
	Endpoint  string `json:"endpoint"`
	Variable  string `json:"variable"`
	RangeFrom int    `json:"range_from"`
	RangeTo   int    `json:"range_to"`
}

func LoadConfig() *Config {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := &Config{}
	decoder.Decode(config)
	return config
}
