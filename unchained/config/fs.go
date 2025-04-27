package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read(file string, params *DynamicParams) (Unchained, error) {
	var cfg Unchained

	f, err := os.Open(file)
	if err != nil {
		return cfg, fmt.Errorf("open config file: %w", err)
	}

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("decode config file: %w", err)
	}
	params.Set(&cfg)

	return cfg, nil
}

func Write(cfg Unchained, file string) error {
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
