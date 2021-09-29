package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadJson_1Count_AllPass(t *testing.T) {

	//read in expected JSON from file
	expectedJsonFilePath := "./test/data/expected/test-result-crypto-hash-1-count-pass.json"
	expectedJsonBytes, err := ioutil.ReadFile(expectedJsonFilePath)
	assert.Nil(t, err)
	assert.NotEmpty(t, expectedJsonBytes)

	//generate actual JSON by sending raw "go test" JSON
	rawJsonFilePath := "./test/data/raw/test-result-crypto-hash-1-count-pass.json"
	testRun := processTestRun(rawJsonFilePath)
	actualJsonBytes, err := json.Marshal(testRun)
	assert.Nil(t, err)
	assert.NotEmpty(t, actualJsonBytes)

	expectedJson := string(expectedJsonBytes)
	actualJson := string(actualJsonBytes)

	//compare the 2 JSONs
	assert.JSONEqf(t, expectedJson, actualJson, "TestRun not equal")
}
