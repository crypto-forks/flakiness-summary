package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

// Raw JSON result step from `go test -json` execution
// There are 4 types of JSON results (specified by Action value) generated for each test run: run, output, pass, fail
// sequence of result steps Action types per test:
// 1. run (once)
// 2. output (one to many)
// 3. pass OR fail (once)

type RawTestStep struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output"`
	Elapsed float32   `json:"Elapsed"`
}

type JobOutput struct {
	CommitSha  string                      `json:"CommitSha"`
	CommitDate time.Time                   `json:"CommitDate"`
	JobStarted time.Time                   `json:"JobStarted"`
	Results    map[string][]TestResult_old `json:"Results"`
}

type TestResult_old struct {
	Result  string  `json:"Result"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
	Test    string  `json:"Test"`
}

type TestResult struct {
	Test    string   `json:"test"`
	Package string   `json:"package"`
	Output  []string `json:"output"`
	Result  string   `json:"result"`
	Elapsed float32  `json:"elapsed"`
}

type PackageResult struct {
	Package string       `json:"package"`
	Result  string       `json:"result"`
	Elapsed float64      `json:"elapsed"`
	Output  []string     `json:"output"`
	Tests   []TestResult `json:"tests"`
}

//converts raw JSON output from "go test" to this struct
type TestRun struct {
	CommitSha      string          `json:"commit_sha"`
	CommitDate     string          `json:"commit_date"`
	JobRunDate     string          `json:"job_run_date"`
	PackageResults []PackageResult `json:"results"` // {
	//PackageResult PackageResult
	// Package       string          `json:"package"`
	// Result        string          `json:"result"`
	// Elapsed       float64         `json:"elapsed"`
	// Output        []string        `json:"output"`
	// Tests         []TestResult    `json:"tests"`
	// Tests   []struct {
	// 	Test    string   `json:"test"`
	// 	Package string   `json:"package"`
	// 	Output  []string `json:"output"`
	// 	Result  string   `json:"result"`
	// 	Elapsed int      `json:"elapsed"`
	// } `json:"tests"`
	//} `json:"results"`
}

type TestResultPackage struct {
	Packages uint
	Tests    uint
}

func processTestRun(rawJsonFilePath string) TestRun {
	var packageResult1 PackageResult
	packageResult1.Elapsed = 0.349
	packageResult1.Output = []string{"PASS\n", "ok  \tgithub.com/onflow/flow-go/crypto/hash\t0.349s\n"}
	packageResult1.Package = "github.com/onflow/flow-go/crypto/hash"
	//packageResult1.Result = "pass"

	f, err := os.Open(rawJsonFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//map of test results - sorted by package name
	//testMap := make(map[string][]TestResult)

	testMap := make(map[string]TestResult)

	//testResult.Results = append(testResult.Results, )
	//var results Result

	s := bufio.NewScanner(f)
	for s.Scan() {
		var rawTestStep RawTestStep
		json.Unmarshal(s.Bytes(), &rawTestStep)

		//most raw test steps will have Test value - only package specific steps won't
		if rawTestStep.Test != "" {

			//check if this test already exists in the map - if so, we add data to it

			//first check if this test exists in the map - create it if doesn't
			_, ok := testMap[rawTestStep.Test]
			if !ok {
				//if it doesn't exist, we create a new TestResult add it to the map
				var newTest TestResult
				newTest.Test = rawTestStep.Test

				//store outputs as a slice of strings - that's how "go test -json" outputs each output string on a separate line
				//for passing tests, there are 2 outputs and for failing tests there are more outputs
				newTest.Output = make([]string, 0)

				testMap[rawTestStep.Test] = newTest
			}

			//second, check if that test exists in the package - create it if it doesn't, update it if does

			testResult := testMap[rawTestStep.Test]

			switch rawTestStep.Action {
			case "run":
				testResult.Package = rawTestStep.Package
				testMap[rawTestStep.Test] = testResult

			case "output":
				testResult.Output = append(testResult.Output, rawTestStep.Output)
				testMap[rawTestStep.Test] = testResult

			case "pass":
				testResult.Result = rawTestStep.Action
				testResult.Elapsed = rawTestStep.Elapsed
				testMap[rawTestStep.Test] = testResult

			case "fail":
				testResult.Result = rawTestStep.Action
				testResult.Elapsed = rawTestStep.Elapsed
				testMap[rawTestStep.Test] = testResult
			}

		} else {
			//package level messages won't have a Test value
			switch rawTestStep.Action {
			case "output":
			case "pass":
				packageResult1.Result = "pass"
			case "fail":
				packageResult1.Result = "fail"
			default:
				panic(fmt.Sprintf("unexpected action: %s", rawTestStep.Action))
			}

			if rawTestStep.Action == "pass" || rawTestStep.Action == "fail" {
				fmt.Print("Test: ")
				fmt.Print(rawTestStep.Test)

				fmt.Print(" Action: ")
				fmt.Print(rawTestStep.Action)

				fmt.Print(" Package: ")
				fmt.Print(rawTestStep.Package)

				fmt.Print(" Elapsed: ")
				fmt.Println(rawTestStep.Elapsed)
			}
		}
	}

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range testMap {
		packageResult1.Tests = append(packageResult1.Tests, v)
	}

	//packageResult1.Tests = []TestResult{testResult1, testResult2, testResult3, testResult4, testResult5, testResult6, testResult7, testResult8, testResult9}

	sort.Slice(packageResult1.Tests, func(i, j int) bool {
		return packageResult1.Tests[i].Test < packageResult1.Tests[j].Test
	})

	var testRun TestRun
	testRun.CommitDate = "Tue Sep 21 18:06:25 2021 -0700"
	testRun.CommitSha = "46baf6c6be29af9c040bc14195e195848598bbae"
	testRun.JobRunDate = "Tue Sep 21 21:06:25 2021 -0700"
	testRun.PackageResults = []PackageResult{packageResult1}

	return testRun
}

func main() {
	commitDate, err := time.Parse(time.RFC3339, os.Getenv("COMMIT_DATE"))
	if err != nil {
		panic(err)
	}

	fmt.Println("commit date: " + commitDate.String())

	jobStarted, err := time.Parse(time.RFC3339, os.Getenv("JOB_STARTED"))
	if err != nil {
		panic(err)
	}

	fmt.Println("job started date: " + jobStarted.String())

	commit := os.Getenv("COMMIT_SHA")
	fmt.Println("commit: " + commit)

	jobOutput := JobOutput{
		CommitSha:  os.Getenv("COMMIT_SHA"),
		CommitDate: commitDate,
		JobStarted: jobStarted,
		Results:    make(map[string][]TestResult_old),
	}

	fmt.Println(jobOutput)

	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	var event Event
	// 	if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
	// 		panic(err)
	// 	}

	// 	if event.Test == nil {
	// 		continue
	// 	}

	// 	switch event.Action {
	// 	case "run":
	// 		// TODO
	// 	case "pass":
	// 		// TODO
	// 	case "fail":
	// 		// TODO
	// 	case "output":
	// 		// TODO
	// 	default:
	// 		panic(fmt.Sprintf("unexpected action: %s", event.Action))
	// 	}

	// }

	// if err := scanner.Err(); err != nil {
	// 	panic(err)
	// }

	// TODO: write results to DB

	// TODO: output a flakiness summary
	// see https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-output-parameter
}
