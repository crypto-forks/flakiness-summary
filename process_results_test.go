package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
		"2Count_Pass": {
			expectedJsonFilePath: "./test/data/expected/test-result-crypto-hash-2-count-pass.json",
			rawJsonFilePath:      "./test/data/raw/test-result-crypto-hash-2-count-pass.json",
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
		expectedTestRun.PackageResults[k].TestMap = make(map[string][]TestResult)
	}

	assert.NoError(t, os.Setenv("COMMIT_DATE", "Tue Sep 21 18:06:25 2021 -0700"))
	assert.NoError(t, os.Setenv("COMMIT_SHA", "46baf6c6be29af9c040bc14195e195848598bbae"))
	assert.NoError(t, os.Setenv("JOB_DATE", "Tue Sep 21 21:06:25 2021 -0700"))

	//simulate generating raw "go test -json" output by loading output from saved file
	actualTestRun := processTestRun(rawJsonFilePath)

	checkTestRuns(t, expectedTestRun, actualTestRun)
}

func checkTestRuns(t *testing.T, expectedTestRun TestRun, actualTestRun TestRun) {
	//it's difficult to determine why 2 test runs aren't equal, so we will check the different sub components of them to see where a potential discrepancy exists
	assert.Equal(t, expectedTestRun.CommitDate, actualTestRun.CommitDate)
	assert.Equal(t, expectedTestRun.CommitSha, actualTestRun.CommitSha)
	assert.Equal(t, expectedTestRun.JobRunDate, actualTestRun.JobRunDate)
	assert.Equal(t, len(expectedTestRun.PackageResults), len(actualTestRun.PackageResults))

	//check each package
	for packageIndex := range expectedTestRun.PackageResults {
		assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Elapsed, actualTestRun.PackageResults[packageIndex].Elapsed)
		assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Package, actualTestRun.PackageResults[packageIndex].Package)
		assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Result, actualTestRun.PackageResults[packageIndex].Result)
		assert.Empty(t, expectedTestRun.PackageResults[packageIndex].TestMap, actualTestRun.PackageResults[packageIndex].TestMap)

		//check outputs of each package result
		assert.Equal(t, len(expectedTestRun.PackageResults[packageIndex].Output), len(actualTestRun.PackageResults[packageIndex].Output))
		for packageOutputIndex := range expectedTestRun.PackageResults[packageIndex].Output {
			assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Output[packageOutputIndex], actualTestRun.PackageResults[packageIndex].Output[packageOutputIndex])
		}

		//check all tests results of each package
		assert.Equal(t, len(expectedTestRun.PackageResults[packageIndex].Tests), len(actualTestRun.PackageResults[packageIndex].Tests))
		for testResultIndex := range expectedTestRun.PackageResults[packageIndex].Tests {
			assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Package, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Package)
			assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Test, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Test)
			assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Elapsed, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Elapsed)
			assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Result, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Result)

			//check all outputs of each test result
			assert.Equal(t, len(expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output), len(actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output))
			for testResultOutputIndex := range expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output {
				assert.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output[testResultOutputIndex], actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output[testResultOutputIndex], fmt.Sprintf("error message; PackageResult[%d] TestResult[%d] Output[%d]", packageIndex, testResultIndex, testResultOutputIndex))
			}
		}
	}
	assert.Equal(t, expectedTestRun, actualTestRun)
}
