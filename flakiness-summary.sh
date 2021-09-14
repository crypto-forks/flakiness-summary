#!/bin/sh -l

# clone the repo
git clone https://github.com/onflow/flow-go.git
cd flow-go

# checkout specified commit
git checkout $1

# setup environment
make install-tools tidy generate-mocks

# run tests 
GO111MODULE=on go test -json -count $NUM_RUNS --tags relic ./... | python process_results.py