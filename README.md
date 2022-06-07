
# ratelimiter

Utility to put in front of a webservice to handle a ratelimit.

## Configuration

With a configuration file :

```
```

Or environment variable :

```
RATELIMIT_TARGET
RATELIMIT_DAEMONPORT
RATELIMIT_LOGLEVEL
RATELIMIT_NUMBER
RATELIMIT_DURATIONINSECONDS
```

cd gowork
CGO_ENABLED=0 go build  github.com/sgaunet/ratelimiter