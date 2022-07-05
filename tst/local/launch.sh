#!/usr/bin/env bash

CWD=$(dirname $0)
cd $CWD/../..

task
./ratelimiter -c tst/local/conf.yaml