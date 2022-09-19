#!/bin/bash

. "$PWD"/devops/utils.sh

# build sonar-v2
cd sonar-lang || exit 1
go build -v ./...
test_exit_code
cd ..