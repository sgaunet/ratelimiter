[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/ratelimiter)](https://goreportcard.com/report/github.com/sgaunet/ratelimiter)
[![GitHub release](https://img.shields.io/github/release/sgaunet/ratelimiter.svg)](https://github.com/sgaunet/ratelimiter/releases/latest)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/ratelimiter/total)
[![Maintainability](https://api.codeclimate.com/v1/badges/6a77e40e92008bab6661/maintainability)](https://codeclimate.com/github/sgaunet/ratelimiter/maintainability)
[![License](https://img.shields.io/github/license/sgaunet/ratelimiter.svg)](LICENSE)

# ratelimiter

Utility to put in front of a webservice to handle a ratelimit. 

Works :

* with http protocol (websocket not tested)
* ratelimit by IP (check X-FORWARDED-FOR header)

## Configuration

With a configuration file :

```
---
logLevel: debug
rateNumber: 1  
rateDurationInSeconds: 1
tagetService: http://localhost:8080
# targetService: http://localhost:5678
daemonPort: 1337
```

Or environment variable :

```
RATELIMIT_TARGET
RATELIMIT_DAEMONPORT
RATELIMIT_LOGLEVEL
RATELIMIT_NUMBER
RATELIMIT_DURATIONINSECONDS
```

# Execution

With a docker-compose (change the latest to the latest version) :


```
version: '3.2'
services:
    echo:
        image: sgaunet/http-echo:latest

    ratelimiter:
      image: sgaunet/ratelimiter:latest
      ports:
        - 1337:1337
      environment:
        - RATELIMIT_TARGET=http://echo:8080
        - RATELIMIT_DAEMONPORT=1337
        - RATELIMIT_LOGLEVEL=debug
        - RATELIMIT_NUMBER=305
        - RATELIMIT_DURATIONINSECONDS=10
```

Or with the single binary, option -c to specify a configuration file.

**Please, do not use latest tag. I'm using it when developping so it could not work properly.**

# Install as a systemd service

```
curl -s https://raw.githubusercontent.com/sgaunet/ratelimiter/main/install.sh | sudo bash
```

**This script has been only tested on amd64 arch.**

# Development

## prerequisites

This project is using :

* golang
* [task for development](https://taskfile.dev/#/)
* docker
* [docker buildx](https://github.com/docker/buildx)
* docker manifest
* [goreleaser](https://goreleaser.com/)
* [vegeta](https://github.com/tsenart/vegeta) : load testing tool
* [pre-commit](https://pre-commit.com/)

There are hooks executed in the precommit stage. Once the project cloned on your disk, please install pre-commit:

```
brew install pre-commit
```

Install tools:

```
task install-prereq
```

And install the hooks:

```
task install-pre-commit
```

If you like to launch manually the pre-commmit hook:

```
task pre-commit
```


## tasks

* compile

```
task 
```

* create docker image

```
task image
```

* test release

```
task snapshot
```

* create release

```
task release
```

## load test

### the binary

```
cd tst/local
docker-compose up -d    # launch the web server that will be behind the ratelimiter
./launch-binary.sh      # execute the binary in front of the webserver
```

Launch the load test with: task load-test

## docker image (env var)

```
cd tst/config-env-vars
docker-compose up -d    # launch the web server and ratelimiter
```

Launch the load test with: task load-test

## docker image (configuration file)

```
cd tst/config-file
docker-compose up -d    # launch the web server and ratelimiter
```

Launch the load test with: task load-test

## Project Disclaimer

This software project is released under the MIT License and was created primarily for fun and testing purposes. While it may offer some interesting functionalities, please note:

* Intended Use
* This project is experimental in nature
* It serves as a playground for ideas and concepts
* The code may not be optimized or production-ready

## Recommendation

If you find the features provided by this project useful or intriguing, we strongly recommend exploring more mature and established solutions for your actual needs. This project is not intended to compete with or replace professional-grade software in its domain.

## Contributions

While we appreciate your interest, please understand that this project may not be actively maintained or developed further. Feel free to fork and experiment with the code as per the MIT License terms.
Thank you for your understanding and enjoy exploring!
