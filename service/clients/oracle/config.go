package oracle

import (
	"gopkg.in/yaml.v3"
	"os"
)

type OracleClientConfig struct {
	AVSID         uint64 `yaml:"avs_id"`
	OracleAddress string `yaml:"oracle_address"`
}

type Config struct {
	OracleClients []OracleClientConfig `yaml:"oracle_clients"`
}

func FillConfig(configPath string) (Config, error) {
	// read in config
	file, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
