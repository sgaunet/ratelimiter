package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ConfigYaml describes the struct of the configuration file of the application
type ConfigYaml struct {
	RateNumber            int    `yaml:"rateNumber"`
	RateDurationInSeconds int    `yaml:"rateDurationInSeconds"`
	TargetService         string `yaml:"targetService"`
	LogLevel              string `yaml:"logLevel"`
	DaemonPort            int    `yaml:"daemonPort"`
}

func (c *ConfigYaml) SetDefaultValue() error {
	if c.DaemonPort == 0 {
		c.DaemonPort = 1337
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.TargetService == "" {
		return errors.New("TargetService not specified.")
	}
	return nil
}

func ReadYamlCnxFile(filename string) (ConfigYaml, error) {
	var config ConfigYaml

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Errorf("Error reading YAML file: %s\n", err)
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logrus.Errorf("Error parsing YAML file: %s\n", err)
		return config, err
	}
	return config, err
}

func initTrace(debugLevel string) *logrus.Logger {
	appLog := logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	appLog.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false,
		FullTimestamp:    false,
		DisableTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	appLog.SetOutput(os.Stdout)

	switch debugLevel {
	case "info":
		appLog.SetLevel(logrus.InfoLevel)
	case "warn":
		appLog.SetLevel(logrus.WarnLevel)
	case "error":
		appLog.SetLevel(logrus.ErrorLevel)
	default:
		appLog.SetLevel(logrus.DebugLevel)
	}
	return appLog
}
