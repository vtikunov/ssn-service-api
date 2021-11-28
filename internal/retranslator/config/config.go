package config

import (
	"os"
	"path/filepath"
	"time"

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

// Retranslator - настройки ретранслятора.
//
// EventChannelSize: размер канала событий между пулами консьюмеров и продьюсеров.
//
// MaxConsumers: определяет максимальное количество воркеров-консьюмеров,
// которые будут запущены конкуретно, при окончании работы какого-либо воркера-консьюмера
// пул консьюмеров будет создавать и запускать следующий в пределах указанного лимита.
//
// ConsumerTimeout: определяет максимальное время работы каждого нового экземпляра воркера-консьюемра,
// по истечении которого ему будет направлена команда Stop.
//
// ConsumerBatchTime: определяет время между обращениями к репозиторию событий.
//
// ConsumerBatchSize: определяет размер пакета событий, получаемый консьюмером из репозитория
// за одно обращение.
//
// MaxProducers: определяет максимальное количество воркеров-продьюсеров,
// которые будут запущены конкуретно, при окончании работы какого-либо воркера-продьюсера
// пул продьюсеров будет создавать и запускать следующий в пределах указанного лимита.
//
// ProducerTimeout: определяет максимальное время работы каждого нового
// экземпляра воркера-продьюсера, по истечении которого ему будет направлена команда Stop.
//
// ProducerMaxWorkers: пределяет максимальное количество вспомогательных воркеров работы продьюсера
// с репозиторием событий eventRepo, которые будут запущены конкуретно.
type Retranslator struct {
	EventChannelSize uint64 `yaml:"eventChannelSize"`

	MaxConsumers      uint64        `yaml:"maxConsumers"`
	ConsumerTimeout   time.Duration `yaml:"consumerTimeout"`
	ConsumerBatchTime time.Duration `yaml:"consumerBatchTime"`
	ConsumerBatchSize uint64        `yaml:"consumerBatchSize"`

	MaxProducers       uint64        `yaml:"maxProducers"`
	ProducerTimeout    time.Duration `yaml:"producerTimeout"`
	ProducerMaxWorkers uint64        `yaml:"producerMaxWorkers"`
}

// Database - contains all parameters database connection.
type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SslMode  string `yaml:"sslmode"`
	Driver   string `yaml:"driver"`
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

// Jaeger - contains all parameters for metrics information.
type Jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

// Kafka - contains all parameters kafka information.
type Kafka struct {
	Topic           string   `yaml:"topic"`
	Brokers         []string `yaml:"brokers"`
	PartitionFactor uint8    `yaml:"partitionFactor"`
	SendRetryMax    uint64   `yaml:"sendRetryMax"`
}

// Status config for service.
type Status struct {
	Port          int    `yaml:"port"`
	Host          string `yaml:"host"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Retranslator Retranslator `yaml:"retranslator"`
	Project      Project      `yaml:"project"`
	Database     Database     `yaml:"database"`
	Metrics      Metrics      `yaml:"metrics"`
	Jaeger       Jaeger       `yaml:"jaeger"`
	Kafka        Kafka        `yaml:"kafka"`
	Status       Status       `yaml:"status"`
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
