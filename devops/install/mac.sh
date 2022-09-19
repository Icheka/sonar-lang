#!/bin/bash

go install github.com/icheka/sonar-lang/sonar-lang@latest

lines="export GOPATH=\$HOME/go\nexport PATH=\$PATH:\$GOPATH/bin"
echo -e "$lines" >> ~/.zshrc