package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func getEnvIntVar(envVar string) (number int) {
	var err error
	varString := os.Getenv(envVar)
	if varString != "" {
		number, err = strconv.Atoi(varString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
	return number
}

func main() {
	var (
		err        error
		configFile string
		cfg        ConfigYaml
	)

	flag.StringVar(&configFile, "c", "", "Config file")
	flag.Parse()

	if configFile == "" {
		cfg.TargetService = os.Getenv("RATELIMIT_TARGET")
		cfg.DaemonPort = getEnvIntVar("RATELIMIT_DAEMONPORT")
		cfg.LogLevel = os.Getenv("RATELIMIT_LOGLEVEL")
		cfg.RateNumber = getEnvIntVar("RATELIMIT_NUMBER")
		cfg.RateDurationInSeconds = getEnvIntVar("RATELIMIT_DURATIONINSECONDS")
	} else {
		cfg, err = ReadYamlCnxFile(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error occured when reading %s : %s\n", configFile, err.Error())
			os.Exit(1)
		}
	}

	cfg.SetDefaultValue()
	log = initTrace(cfg.LogLevel)
	app, err := NewApp(cfg)
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(1)
	}
	app.LaunchWebServer()
}
