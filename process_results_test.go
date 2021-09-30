package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadJson_1Count_AllPass(t *testing.T) {

	var expectedTestRun TestRun
	//read in expected JSON from file
	expectedJsonFilePath := "./test/data/expected/test-result-crypto-hash-1-count-pass.json"
	expectedJsonBytes, err := ioutil.ReadFile(expectedJsonFilePath)
	assert.Nil(t, err)
	assert.NotEmpty(t, expectedJsonBytes)

	json.Unmarshal(expectedJsonBytes, &expectedTestRun)
	//order Tests alphabetically so it matches actual output - otherwise, equality check will fail
	sort.Slice(expectedTestRun.PackageResults[0].Tests, func(i, j int) bool {
		return expectedTestRun.PackageResults[0].Tests[i].Test < expectedTestRun.PackageResults[0].Tests[j].Test
	})

	//simulate generating raw "go test -json" output by loading output from saved file
	rawJsonFilePath := "./test/data/raw/test-result-crypto-hash-1-count-pass.json"
	actualTestRun := processTestRun(rawJsonFilePath)
	assert.Nil(t, err)

	assert.Equal(t, expectedTestRun, actualTestRun)
}
