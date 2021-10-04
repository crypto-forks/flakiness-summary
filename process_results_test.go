package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ProcessTestRunData struct {
	expectedJsonFilePath string
	rawJsonFilePath      string
}

//data driven table test
func TestProcessTestRun(t *testing.T) {
	processTestMap2 := map[string]ProcessTestRunData{
		"1Count_AllPass": {
			expectedJsonFilePath: "./test/data/expected/test-result-crypto-hash-1-count-pass.json",
			rawJsonFilePath:      "./test/data/raw/test-result-crypto-hash-1-count-pass.json",
		},
		"1Count_Fail": {
			expectedJsonFilePath: "./test/data/expected/test-result-crypto-hash-1-count-fail.json",
			rawJsonFilePath:      "./test/data/raw/test-result-crypto-hash-1-count-fail.json",
		},
	}

	for k, pt := range processTestMap2 {
		t.Run(k, func(t *testing.T) {
			runProcessTestRun(t, pt.expectedJsonFilePath, pt.rawJsonFilePath)
		})
	}
}

//HELPERS - UTILITIES

func runProcessTestRun(t *testing.T, expectedJsonFilePath string, rawJsonFilePath string) {
	var expectedTestRun TestRun
	//read in expected JSON from file
	expectedJsonBytes, err := ioutil.ReadFile(expectedJsonFilePath)
	assert.Nil(t, err)
	assert.NotEmpty(t, expectedJsonBytes)

	json.Unmarshal(expectedJsonBytes, &expectedTestRun)

	//sort all package alphabetically
	sort.Slice(expectedTestRun.PackageResults, func(i, j int) bool {
		return expectedTestRun.PackageResults[i].Package < expectedTestRun.PackageResults[j].Package
	})

	//sort all tests alphabetically within each package - otherwise, equality check will fail
	for k := range expectedTestRun.PackageResults {
		sort.Slice(expectedTestRun.PackageResults[k].Tests, func(i, j int) bool {
			return expectedTestRun.PackageResults[k].Tests[i].Test < expectedTestRun.PackageResults[k].Tests[j].Test
		})

		//init TestMap to empty - otherwise get comparison failure because would be nil
		expectedTestRun.PackageResults[k].TestMap = make(map[string]TestResult)
	}

	//simulate generating raw "go test -json" output by loading output from saved file
	actualTestRun := processTestRun(rawJsonFilePath)
	assert.Nil(t, err)

	assert.Equal(t, expectedTestRun, actualTestRun)
}
