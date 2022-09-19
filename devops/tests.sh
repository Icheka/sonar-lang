#!/bin/bash

. "$PWD"/devops/utils.sh

# test sonar-lang
cd sonar-lang || exit 1
go test -v ./...
test_exit_code
cd ..

# test language-server
cd language-server || exit 1
go test -v ./...
test_exit_code
cd ..

# test code-editor
cd code-editor || exit 1
export CI=true
npm run test
test_exit_code
cd ..