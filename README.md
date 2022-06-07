
# ratelimiter

Utility to put in front of a webservice to handle a ratelimit. Works only with http protocol actually.

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

With a docker-compose :


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