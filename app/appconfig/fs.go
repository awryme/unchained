package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadUnchained(file string, params *UnchainedRuntimeParams) (Unchained, error) {
	var cfg Unchained

	f, err := os.Open(file)
	if err != nil {
		return cfg, fmt.Errorf("open config file: %w", err)
	}

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("decode config file: %w", err)
	}

	cfg.setRuntimeParams(params)
	return cfg, nil
}

func WriteUnchained(cfg Unchained, file string) error {
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

func ReadUnchainedWorker(file string, params *UnchainedWorkerRuntimeParams) (UnchainedWorker, error) {
	var cfg UnchainedWorker

	f, err := os.Open(file)
	if err != nil {
		return cfg, fmt.Errorf("open config file: %w", err)
	}

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("decode config file: %w", err)
	}

	cfg.setRuntimeParams(params)
	return cfg, nil
}

func WriteUnchainedWorker(cfg UnchainedWorker, file string) error {
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
