#!/bin/sh -l

set -e

export PATH="/usr/local/go/bin:/usr/bin:$PATH"
export GOPATH=$(go env GOPATH)
export GOBIN=$GOPATH/bin

# clone the repo
git clone https://github.com/onflow/flow-go.git
cd flow-go

# checkout specified commit
if [ -z "$1" ]
then
    git checkout master
else
    git checkout $1
fi

export COMMIT_SHA=$(git rev-parse HEAD)
export COMMIT_TIME=$(git show --no-patch --no-notes --pretty='%ct' $COMMIT_SHA)

# setup environment
make install-tools tidy generate-mocks

# run tests 
export NUM_RUNS=10
GO111MODULE=on go test -json -count $NUM_RUNS --tags relic ./... | python3 process_results.py