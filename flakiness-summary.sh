#!/usr/bin/env bash

set -e

export JOB_STARTED=$(date -Iseconds)

case $1 in
    unit|crypto-unit|integration-unit|integration)
        category=$1
    ;;
    *)
        echo "Valid test category must be provided."
        exit 1
    ;;
esac

process_results="go run $(realpath ./process_results.go)"

# clone the repo
git clone https://github.com/onflow/flow-go.git
cd flow-go

# checkout specified commit
if [ -z "$2" ]
then
    git checkout master
else
    git checkout $2
fi

export COMMIT_SHA=$(git rev-parse HEAD)
export COMMIT_DATE=$(git show --no-patch --no-notes --pretty='%cI' $COMMIT_SHA)

export GOPATH=$(/usr/local/go/bin/go env GOPATH)
export GOBIN=$GOPATH/bin
export PATH="/usr/local/go/bin:$GOBIN:$PATH"

make crypto/relic/build

NUM_RUNS=10

case $category in
    unit)
        cd $GOPATH
        GO111MODULE=on go get github.com/vektra/mockery/cmd/mockery@v1.1.2
        GO111MODULE=on go get github.com/golang/mock/mockgen@v1.3.1
        cd -
        make generate-mocks
        GO111MODULE=on go test -json -count $NUM_RUNS --tags relic ./... | $process_results
    ;;
    crypto-unit)
        cd ./crypto
        GO111MODULE=on go test -json -count $NUM_RUNS --tags relic ./... | $process_results
    ;;
    integration-unit)
        cd ./integration
        GO111MODULE=on go test -json -count $NUM_RUNS --tags relic `go list ./... | grep -v -e integration/tests -e integration/benchmark` | $process_results
    ;;
    integration)
        make docker-build-flow
        cd ./integration/tests
        GO111MODULE=on go test -json -count $NUM_RUNS --tags relic ./... | $process_results
    ;;
esac

