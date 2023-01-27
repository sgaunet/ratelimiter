#!/usr/bin/env bash

CWD=$(dirname "$0")
cd "$CWD/../.." || exit 1

task   # commpile project
./ratelimiter -c tst/local/conf.yaml
