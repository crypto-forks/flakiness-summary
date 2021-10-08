package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

//data driven table test
func TestProcessTestRun(t *testing.T) {
	testDataMap := map[string]string{
		"1 count all pass":               "test-result-crypto-hash-1-count-pass.json",
		"1 count 1 fail the rest pass":   "test-result-crypto-hash-1-count-fail.json",
		"1 count 2 skiped the rest pass": "test-result-crypto-hash-1-count-skip-pass.json",
		"2 count all pass":               "test-result-crypto-hash-2-count-pass.json",
		"10 count all pass":              "test-result-crypto-hash-10-count-pass.json",
		"10 count some failures":         "test-result-crypto-hash-10-count-fail.json",
	}

	for k, testJsonData := range testDataMap {
		t.Run(k, func(t *testing.T) {
			runProcessTestRun(t, testJsonData)
		})
	}
}

//HELPERS - UTILITIES

func runProcessTestRun(t *testing.T, jsonExpectedActualFile string) {
	const expectedJsonFilePath = "./test/data/expected/"
	const rawJsonFilePath = "./test/data/raw/"

	var expectedTestRun TestRun
	//read in expected JSON from file
	expectedJsonBytes, err := ioutil.ReadFile(expectedJsonFilePath + jsonExpectedActualFile)
	require.Nil(t, err)
	require.NotEmpty(t, expectedJsonBytes)

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

	require.NoError(t, os.Setenv("COMMIT_DATE", "Tue Sep 21 18:06:25 2021 -0700"))
	require.NoError(t, os.Setenv("COMMIT_SHA", "46baf6c6be29af9c040bc14195e195848598bbae"))
	require.NoError(t, os.Setenv("JOB_DATE", "Tue Sep 21 21:06:25 2021 -0700"))

	//simulate generating raw "go test -json" output by loading output from saved file
	actualTestRun := processTestRun(rawJsonFilePath + jsonExpectedActualFile)

	checkTestRuns(t, expectedTestRun, actualTestRun)
}

func checkTestRuns(t *testing.T, expectedTestRun TestRun, actualTestRun TestRun) {
	//it's difficult to determine why 2 test runs aren't equal, so we will check the different sub components of them to see where a potential discrepancy exists
	require.Equal(t, expectedTestRun.CommitDate, actualTestRun.CommitDate)
	require.Equal(t, expectedTestRun.CommitSha, actualTestRun.CommitSha)
	require.Equal(t, expectedTestRun.JobRunDate, actualTestRun.JobRunDate)
	require.Equal(t, len(expectedTestRun.PackageResults), len(actualTestRun.PackageResults))

	//check each package
	for packageIndex := range expectedTestRun.PackageResults {
		require.Equal(t, expectedTestRun.PackageResults[packageIndex].Elapsed, actualTestRun.PackageResults[packageIndex].Elapsed)
		require.Equal(t, expectedTestRun.PackageResults[packageIndex].Package, actualTestRun.PackageResults[packageIndex].Package)
		require.Equal(t, expectedTestRun.PackageResults[packageIndex].Result, actualTestRun.PackageResults[packageIndex].Result)
		require.Empty(t, expectedTestRun.PackageResults[packageIndex].TestMap, actualTestRun.PackageResults[packageIndex].TestMap)

		//check outputs of each package result
		require.Equal(t, len(expectedTestRun.PackageResults[packageIndex].Output), len(actualTestRun.PackageResults[packageIndex].Output))
		for packageOutputIndex := range expectedTestRun.PackageResults[packageIndex].Output {
			require.Equal(t, expectedTestRun.PackageResults[packageIndex].Output[packageOutputIndex], actualTestRun.PackageResults[packageIndex].Output[packageOutputIndex])
		}

		//check all tests results of each package
		require.Equal(t, len(expectedTestRun.PackageResults[packageIndex].Tests), len(actualTestRun.PackageResults[packageIndex].Tests))
		for testResultIndex := range expectedTestRun.PackageResults[packageIndex].Tests {

			//check all outputs of each test result
			require.Equal(t, len(expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output), len(actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output), fmt.Sprintf("TestResult[%d].Test: %s", testResultIndex, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Test))
			for testResultOutputIndex := range expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output {
				require.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output[testResultOutputIndex], actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Output[testResultOutputIndex], fmt.Sprintf("PackageResult[%d] TestResult[%d] Output[%d]", packageIndex, testResultIndex, testResultOutputIndex))
			}

			require.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Package, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Package)
			require.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Test, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Test)
			require.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Elapsed, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Elapsed, fmt.Sprintf("TestResult[%d].Test: %s", testResultIndex, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Test))
			require.Equal(t, expectedTestRun.PackageResults[packageIndex].Tests[testResultIndex].Result, actualTestRun.PackageResults[packageIndex].Tests[testResultIndex].Result)
		}
	}
	require.Equal(t, expectedTestRun, actualTestRun)
}
