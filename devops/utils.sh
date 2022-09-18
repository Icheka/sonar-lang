#!/bin/bash

test_exit_code() {
    if [ $? -eq 1 ]
    then
        exit 1
    fi
}