package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

// GetConfigInstance returns service config
func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Project - contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string
	CommitHash  string
}

// Metrics - contains all parameters for metrics information.
type Metrics struct {
	Port      int    `yaml:"port"`
	Host      string `yaml:"host"`
	Path      string `yaml:"path"`
	Subsystem string `yaml:"subsystem"`
}

// Jaeger - contains all parameters for jaeger information.
type Jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

// Status config for service.
type Status struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

// Bot - contains all parameters for telegram.
type Bot struct {
	Timeout                 uint   `yaml:"timeout"`
	ListPerPage             uint64 `yaml:"listPerPage"`
	WriteServiceServiceAddr string `yaml:"writeServiceServiceAddr"`
	ReadServiceServiceAddr  string `yaml:"readServiceServiceAddr"`
	ServicesCallRetries     uint64 `yaml:"servicesCallRetries"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project Project `yaml:"project"`
	Metrics Metrics `yaml:"metrics"`
	Jaeger  Jaeger  `yaml:"jaeger"`
	Status  Status  `yaml:"status"`
	Bot     Bot     `yaml:"bot"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
