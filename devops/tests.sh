#!/bin/bash

. "$PWD"/devops/utils.sh

# test sonar-v2
cd sonar-lang || exit 1
go test -v ./...
test_exit_code
cd ..