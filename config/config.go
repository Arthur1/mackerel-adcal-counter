package config

import (
	"errors"
	"net/url"
)

type Config struct {
	Url            string
	Service        string
	MetricName     string
	MackerelApiKey string
}

func NewConfig(rawurl string, service string, metricName string, mackerelApiKey string) (*Config, error) {
	if err := isValidUrl(rawurl); err != nil {
		return nil, err
	}

	if err := isValidService(service); err != nil {
		return nil, err
	}

	if err := isValidMetricName(metricName); err != nil {
		return nil, err
	}

	config := &Config{
		Url:            rawurl,
		Service:        service,
		MetricName:     metricName,
		MackerelApiKey: mackerelApiKey,
	}
	return config, nil
}

func isValidUrl(rawurl string) error {
	_, err := url.ParseRequestURI(rawurl)
	return err
}

func isValidService(service string) error {
	if service == "" {
		return errors.New("service is required")
	}
	// TODO: 正規表現を用いた Validation
	return nil
}

func isValidMetricName(metricName string) error {
	if metricName == "" {
		return errors.New("metricName is required")
	}
	// TODO: 正規表現を用いた Validation
	return nil
}

func (c *Config) MackerelApiKeyIsSet() bool {
	return c.MackerelApiKey != ""
}
