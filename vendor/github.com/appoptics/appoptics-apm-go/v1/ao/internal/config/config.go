// Copyright (C) 2017 Librato, Inc. All rights reserved.

// Package config is responsible for loading the configuration from various
// sources, e.g., environment variables, configuration files and user input.
// It also accepts dynamic settings from the collector server.
package config

import (
	"github.com/appoptics/appoptics-apm-go/v1/ao/internal/log"
)

// The default values for environment variables
const (
	defaultGRPCCollector      = "collector.appoptics.com:443"
	defaultServiceKey         = ""
	defaultTrustedPath        = ""
	defaultCollectorUDP       = "127.0.0.1:7831"
	defaultReporter           = "ssl"
	defaultTracingMode        = "always"
	defaultPrependDomain      = false
	defaultHostnameAlias      = ""
	defaultInsecureSkipVerify = false
	defaultHistogramPrecision = 2
)

// The environment variables
const (
	envAppOpticsCollector           = "APPOPTICS_COLLECTOR"
	envAppOpticsServiceKey          = "APPOPTICS_SERVICE_KEY"
	envAppOpticsTrustedPath         = "APPOPTICS_TRUSTEDPATH"
	envAppOpticsCollectorUDP        = "APPOPTICS_COLLECTOR_UDP"
	envAppOpticsReporter            = "APPOPTICS_REPORTER"
	envAppOpticsTracingMode         = "APPOPTICS_TRACING_MODE"
	envAppOpticsPrependDomain       = "APPOPTICS_PREPEND_DOMAIN"
	envAppOpticsHostnameAlias       = "APPOPTICS_HOSTNAME_ALIAS"
	envAppOpticsInsecureSkipVerify  = "APPOPTICS_INSECURE_SKIP_VERIFY"
	envAppOpticsHistogramPrecision  = "APPOPTICS_HISTOGRAM_PRECISION"
	envAppOpticsEventsFlushInterval = "APPOPTICS_EVENTS_FLUSH_INTERVAL"
	envAppOpticsEventsBatchSize     = "APPOPTICS_EVENTS_BATCHSIZE"
)

// The environment variables, validators and converters. This map is not
// concurrent-safe and should not be modified after initialization.
var envs = map[string]Env{
	"Collector": {
		name:     envAppOpticsCollector,
		optional: true,
		validate: IsValidHost,
		convert:  ToHost,
		mask:     nil,
	},
	"ServiceKey": {
		name:     envAppOpticsServiceKey,
		optional: false,
		validate: IsValidServiceKey,
		convert:  ToServiceKey,
		mask:     maskServiceKey,
	},
	"TrustedPath": {
		name:     envAppOpticsTrustedPath,
		optional: true,
		validate: IsValidFileString,
		convert:  ToFileString,
		mask:     nil,
	},
	"CollectorUDP": {
		name:     envAppOpticsCollectorUDP,
		optional: true,
		validate: IsValidHost,
		convert:  ToHost,
		mask:     nil,
	},
	"ReporterType": {
		name:     envAppOpticsReporter,
		optional: true,
		validate: IsValidReporterType,
		convert:  ToReporterType,
		mask:     nil,
	},
	"TracingMode": {
		name:     envAppOpticsTracingMode,
		optional: true,
		validate: IsValidTracingMode,
		convert:  ToTracingMode,
		mask:     nil,
	},
	"PrependDomain": {
		name:     envAppOpticsPrependDomain,
		optional: true,
		validate: IsValidBool,
		convert:  ToBool,
		mask:     nil,
	},
	"HostAlias": {
		name:     envAppOpticsHostnameAlias,
		optional: true,
		validate: IsValidHostnameAlias,
		convert:  ToHostnameAlias,
		mask:     nil,
	},
	"SkipVerify": {
		name:     envAppOpticsInsecureSkipVerify,
		optional: true,
		validate: IsValidBool,
		convert:  ToBool,
		mask:     nil,
	},
	"Precision": {
		name:     envAppOpticsHistogramPrecision,
		optional: true,
		validate: IsValidInteger,
		convert:  ToInteger,
		mask:     nil,
	},
	"EventsFlushInterval": {
		name:     envAppOpticsEventsFlushInterval,
		optional: true,
		validate: IsValidInteger,
		convert:  ToInt64,
		mask:     nil,
	},
	"EventsBatchSize": {
		name:     envAppOpticsEventsBatchSize,
		optional: true,
		validate: IsValidInteger,
		convert:  ToInt64,
		mask:     nil,
	},
}

// Config is the struct to define the agent configuration. The configuration
// options in this struct (excluding those from ReporterOptions) are not
// intended for dynamically updating.
type Config struct {
	// Collector defines the host and port of the AppOptics collector
	Collector string `yaml:"CollectorHost" json:"CollectorHost"`

	// ServiceKey defines the service key and service name
	ServiceKey string

	// The file path of the cert file for gRPC connection
	TrustedPath string

	// The host and port of the UDP collector
	CollectorUDP string `yaml:"CollectorHostUDP" json:"CollectorHostUDP"`

	// The reporter type, ssl or udp
	ReporterType string

	// The tracing mode
	TracingMode string

	// Whether the domain should be prepended to the transaction name.
	PrependDomain bool

	// The alias of the hostname
	HostAlias string `yaml:"HostnameAlias" json:"HostnameAlias"`

	// Whether to skip verification of hostname
	SkipVerify bool `yaml:"InsecureSkipVerify" json:"InsecureSkipVerify"`

	// The precision of the histogram
	Precision int `yaml:"HistogramPrecision" json:"HistogramPrecision"`

	// The reporter options
	Reporter *ReporterOptions `yaml:"ReporterOptions" json:"ReporterOptions"`
}

// Option is a function type that accepts a Config pointer and
// applies the configuration option it defines.
type Option func(c *Config)

// WithCollector defines a Config option for collector address.
func WithCollector(collector string) Option {
	return func(c *Config) {
		c.Collector = collector
	}
}

// WithServiceKey defines a Config option for the service key.
func WithServiceKey(key string) Option {
	return func(c *Config) {
		c.ServiceKey = key
	}
}

// NewConfig initializes a ReporterOptions object and override default values
// with options provided as arguments. It may print errors if there are invalid
// values in the configuration file or the environment variables.
//
// It returns a config with best-effort, e.g., fall back to default values for
// invalid environment variables. The reporter will be the final decision maker
// on whether to start up.
func NewConfig(opts ...Option) *Config {
	c := newConfig()
	c.RefreshConfig(opts...)
	return c
}

// RefreshConfig loads the customized settings and merge with default values
func (c *Config) RefreshConfig(opts ...Option) {
	c.LoadConfigFile("TODO")
	c.LoadEnvs()

	for _, opt := range opts {
		opt(c)
	}
}

func newConfig() *Config {
	return &Config{
		Collector:     defaultGRPCCollector,
		ServiceKey:    defaultServiceKey,
		TrustedPath:   defaultTrustedPath,
		CollectorUDP:  defaultCollectorUDP,
		ReporterType:  defaultReporter,
		TracingMode:   defaultTracingMode,
		PrependDomain: defaultPrependDomain,
		HostAlias:     defaultHostnameAlias,
		SkipVerify:    defaultInsecureSkipVerify,
		Precision:     defaultHistogramPrecision,
		Reporter:      defaultReporterOptions(),
	}
}

// LoadEnvs loads environment variable values and update the Config object.
func (c *Config) LoadEnvs() {
	// TODO: reflect?
	c.Collector = envs["Collector"].LoadString(c.Collector)
	c.ServiceKey = envs["ServiceKey"].LoadString(c.ServiceKey)

	c.TrustedPath = envs["TrustedPath"].LoadString(c.TrustedPath)
	c.CollectorUDP = envs["CollectorUDP"].LoadString(c.CollectorUDP)
	c.ReporterType = envs["ReporterType"].LoadString(c.ReporterType)
	c.TracingMode = envs["TracingMode"].LoadString(c.TracingMode)

	c.PrependDomain = envs["PrependDomain"].LoadBool(c.PrependDomain)
	c.HostAlias = envs["HostAlias"].LoadString(c.HostAlias)
	c.SkipVerify = envs["SkipVerify"].LoadBool(c.SkipVerify)

	c.Precision = envs["Precision"].LoadInt(c.Precision)

	c.Reporter.LoadEnvs()
}

// LoadConfigFile loads from the config file
func (c *Config) LoadConfigFile(path string) error {
	log.Debug("Loading from config file is not implemented.")
	return nil
}

// GetCollector returns the collector address
func (c *Config) GetCollector() string {
	return c.Collector
}

// GetServiceKey returns the service key
func (c *Config) GetServiceKey() string {
	return c.ServiceKey
}

// GetTrustedPath returns the file path of the cert file
func (c *Config) GetTrustedPath() string {
	return c.TrustedPath
}

// GetReporterType returns the reporter type
func (c *Config) GetReporterType() string {
	return c.ReporterType
}

// GetCollectorUDP returns the UDP collector host
func (c *Config) GetCollectorUDP() string {
	return c.CollectorUDP
}

// GetTracingMode returns the UDP collector host
func (c *Config) GetTracingMode() string {
	return c.TracingMode
}

// GetPrependDomain returns the UDP collector host
func (c *Config) GetPrependDomain() bool {
	return c.PrependDomain
}

// GetHostAlias returns the UDP collector host
func (c *Config) GetHostAlias() string {
	return c.HostAlias
}

// GetSkipVerify returns the UDP collector host
func (c *Config) GetSkipVerify() bool {
	return c.SkipVerify
}

// GetPrecision returns the UDP collector host
func (c *Config) GetPrecision() int {
	return c.Precision
}

// GetReporter returns the reporter options struct
func (c *Config) GetReporter() *ReporterOptions {
	return c.Reporter
}
