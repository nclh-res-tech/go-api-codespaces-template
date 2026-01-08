package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"{{MODULE_PATH}}/common/errors"
	"{{MODULE_PATH}}/common/observability"
	corehttp "{{MODULE_PATH}}/internal/core/httpserver"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	// ErrInvalidEnvironment is returned when the API_ENVIRONMENT is not set or invalid.
	ErrInvalidEnvironment = errors.Error("API_ENVIRONMENT is not set")
	// ErrValidation is returned when the configuration is invalid.
	ErrValidation = errors.Error("invalid configuration")
	// ErrRead is returned when the configuration file cannot be read.
	ErrRead = errors.Error("failed to read file")
	// ErrUnmarshal is returned when the configuration file cannot be unmarshalled.
	ErrUnmarshal = errors.Error("failed to unmarshal file")
)

var (
	baseConfigPath = "config/config.yaml"
	envConfigPath  = "config/config-%s.yaml"
)

// Settings holds backend connectivity options (REST/SOAP).
type Settings struct {
	BackendURL           string            `koanf:"url" json:"backend_url" yaml:"backend_url"`
	BackendMode          string            `koanf:"mode" json:"backend_mode,omitempty" yaml:"backend_mode,omitempty"`
	SOAPAction           string            `koanf:"soap_action" json:"soap_action,omitempty" yaml:"soap_action,omitempty"`
	SOAPEnvelopeTemplate string            `koanf:"soap_envelope_template" json:"soap_envelope_template,omitempty" yaml:"soap_envelope_template,omitempty"`
	BackendURLOverrides  map[string]string `koanf:"url_overrides" json:"backend_url_overrides,omitempty" yaml:"backend_url_overrides,omitempty"`
}

// Config represents the configuration of our application.
type Config struct {
	HTTP      corehttp.Config      `koanf:"http"`
	Auth      AuthConfig           `koanf:"auth"`
	Telemetry observability.Config `koanf:"telemetry"`
	Backend   Settings             `koanf:"backend"`
	APIEnv    string               `koanf:"api_environment" validate:"required"`
}

// AuthConfig controls HTTP authentication middleware.
type AuthConfig struct {
	Enabled bool `koanf:"enabled"`

	Cert struct {
		Enabled          bool     `koanf:"enabled"`
		CertFingerprints []string `koanf:"cert_fingerprints"`
		CertCommonNames  []string `koanf:"cert_common_names"`
	} `koanf:"cert"`
}

// Load loads the configuration from yaml files and environment variables using koanf.
func Load(ctx context.Context) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(baseConfigPath), yaml.Parser()); err != nil {
		return nil, ErrRead.Wrap(err)
	}

	environ := os.Getenv("API_ENVIRONMENT")

	if environ != "" {
		envConfig := fmt.Sprintf(envConfigPath, environ)
		if _, err := os.Stat(envConfig); err == nil {
			if err := k.Load(file.Provider(envConfig), yaml.Parser()); err != nil {
				return nil, ErrUnmarshal.Wrap(err)
			}
		}
	}

	if err := k.Load(env.Provider("APP_", ".", func(key string) string {
		key = strings.TrimPrefix(key, "APP_")
		key = strings.ReplaceAll(strings.ToLower(key), "__", ".")
		return key
	}), nil); err != nil {
		return nil, ErrUnmarshal.Wrap(err)
	}

	config := &Config{}
	if err := k.Unmarshal("", config); err != nil {
		return nil, ErrUnmarshal.Wrap(err)
	}

	if config.APIEnv == "" && environ != "" {
		config.APIEnv = environ
	}

	if config.APIEnv == "" {
		return nil, ErrInvalidEnvironment
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, ErrValidation.Wrap(err)
	}

	fmt.Printf("loaded config environment: %s\n", config.APIEnv)

	return config, nil
}

// LoadSettings returns only the backend settings block.
func LoadSettings() (Settings, error) {
	cfg, err := Load(context.Background())
	if err != nil {
		return Settings{}, err
	}
	return cfg.Backend, nil
}
