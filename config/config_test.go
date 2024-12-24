package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	fileContent := `{
		"cells": [
			{
				"name": "Cell1",
				"endpoint": "http://example.com",
				"variable": "var1",
				"range_from": 1,
				"range_to": 100
			}
		]
	}`
	file, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(fileContent)); err != nil {
		t.Fatal(err)
	}
	file.Close()

	// Override the config file path for the test
	oldConfigPath := configFilePath
	configFilePath = file.Name()
	defer func() { configFilePath = oldConfigPath }()

	config := LoadConfig()
	if len(config.Cells) != 1 {
		t.Errorf("Expected 1 cell, got %d", len(config.Cells))
	}
	if config.Cells[0].Name != "Cell1" {
		t.Errorf("Expected cell name 'Cell1', got %s", config.Cells[0].Name)
	}
}
