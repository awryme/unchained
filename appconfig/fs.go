package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read(file string, params *RuntimeParams) (Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return Config{}, fmt.Errorf("open config file: %w", err)
	}
	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decode config file: %w", err)
	}
	setRuntimeParams(&cfg, params)
	return cfg, nil
}

func Write(cfg Config, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create config file: %w", err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(cfg)
	if err != nil {
		return fmt.Errorf("encode config file: %w", err)
	}
	return nil
}
