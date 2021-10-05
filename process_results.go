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

//models single line from "go test -json" output
type RawTestStep struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output"`
	Elapsed float32   `json:"Elapsed"`
}

//models full summary of a test run from "go test -json"
type TestRun struct {
	CommitSha      string          `json:"commit_sha"`
	CommitDate     string          `json:"commit_date"`
	JobRunDate     string          `json:"job_run_date"`
	PackageResults []PackageResult `json:"results"` // {
}

//models test result of an entire package which can have multiple tests
type PackageResult struct {
	Package string       `json:"package"`
	Result  string       `json:"result"`
	Elapsed float32      `json:"elapsed"`
	Output  []string     `json:"output"`
	Tests   []TestResult `json:"tests"`
	TestMap map[string][]TestResult
}

//models result of a single test that's part of a larger package result
type TestResult struct {
	Test    string   `json:"test"`
	Package string   `json:"package"`
	Output  []string `json:"output"`
	Result  string   `json:"result"`
	Elapsed float32  `json:"elapsed"`
}

//models github CI test job
type JobOutput struct {
	CommitSha  string                `json:"CommitSha"`
	CommitDate time.Time             `json:"CommitDate"`
	JobStarted time.Time             `json:"JobStarted"`
	Results    map[string]TestResult `json:"Results"`
}

func processTestRun(rawJsonFilePath string) TestRun {
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

	//map of package results
	packageMap := make(map[string]PackageResult)

	s := bufio.NewScanner(f)
	for s.Scan() {
		var rawTestStep RawTestStep
		json.Unmarshal(s.Bytes(), &rawTestStep)

		//check if package result exists to hold test results
		packageResult, ok := packageMap[rawTestStep.Package]
		if !ok {
			//if package doesn't exist, create new package result and add it to map
			var newPackageResult PackageResult
			newPackageResult.Package = rawTestStep.Package

			//store outputs as a slice of strings - that's how "go test -json" outputs each output string on a separate line
			//there are usually 2 or more outputs for a package
			newPackageResult.Output = make([]string, 0)

			//package result will hold map of test results
			newPackageResult.TestMap = make(map[string][]TestResult)

			packageMap[rawTestStep.Package] = newPackageResult
			packageResult = newPackageResult
		}

		//most raw test steps will have Test value - only package specific steps won't
		if rawTestStep.Test != "" {
			//subsequent raw json outputs will have different data about the test - whether it passed/failed, what the test output was, etc
			switch rawTestStep.Action {

			// Raw JSON result step from `go test -json` execution
			// There are 4 types of JSON results (specified by Action value) generated for each test run: run, output, pass, fail
			// sequence of result steps Action types per test:
			// 1. run (once)
			// 2. output (one to many)
			// 3. pass OR fail (once)

			case "run":
				//testResults, ok := packageResult.TestMap[rawTestStep.Test]

				var newTestResult TestResult
				newTestResult.Test = rawTestStep.Test

				//store outputs as a slice of strings - that's how "go test -json" outputs each output string on a separate line
				//for passing tests, there are usually 2 outputs for a passing test and more outputs for a failing test
				newTestResult.Output = make([]string, 0)

				//if test result doesn't exist, create a new test result add it to the test result slice
				if !ok {

					newTestResults := []TestResult{newTestResult}

					packageResult.TestMap[rawTestStep.Test] = newTestResults

					//testResults = newTestResults
				} else {
					//test result exists but it's a new count / run - append to test result slice
					packageResult.TestMap[rawTestStep.Test] = append(packageResult.TestMap[rawTestStep.Test], newTestResult)
				}
				//lastIndex := len(testResults) - 1
				lastIndex := len(packageResult.TestMap[rawTestStep.Test]) - 1
				//testResults[lastIndex].Package = rawTestStep.Package
				packageResult.TestMap[rawTestStep.Test][lastIndex].Package = rawTestStep.Package

			case "output":
				testResults, ok := packageResult.TestMap[rawTestStep.Test]
				if !ok {
					panic(fmt.Sprintf("no test result for test %s", rawTestStep.Test))
				}
				lastIndex := len(packageResult.TestMap[rawTestStep.Test]) - 1
				//testResults[lastIndex].Output = append(testResults[lastIndex].Output, rawTestStep.Output)
				packageResult.TestMap[rawTestStep.Test][lastIndex].Output = append(testResults[lastIndex].Output, rawTestStep.Output)

			case "pass":
				// testResults, ok := packageResult.TestMap[rawTestStep.Test]
				// if !ok {
				// 	panic(fmt.Sprintf("no test result for test %s", rawTestStep.Test))
				// }
				lastIndex := len(packageResult.TestMap[rawTestStep.Test]) - 1
				packageResult.TestMap[rawTestStep.Test][lastIndex].Result = rawTestStep.Action
				packageResult.TestMap[rawTestStep.Test][lastIndex].Elapsed = rawTestStep.Elapsed
				// testResults[lastIndex].Result = rawTestStep.Action
				// testResults[lastIndex].Elapsed = rawTestStep.Elapsed

			case "fail":
				// testResults, ok := packageResult.TestMap[rawTestStep.Test]
				// if !ok {
				// 	panic(fmt.Sprintf("no test result for test %s", rawTestStep.Test))
				// }
				lastIndex := len(packageResult.TestMap[rawTestStep.Test]) - 1
				packageResult.TestMap[rawTestStep.Test][lastIndex].Result = rawTestStep.Action
				packageResult.TestMap[rawTestStep.Test][lastIndex].Elapsed = rawTestStep.Elapsed

				// testResults[lastIndex].Result = rawTestStep.Action
				// testResults[lastIndex].Elapsed = rawTestStep.Elapsed
			default:
				panic(fmt.Sprintf("unexpected action: %s", rawTestStep.Action))
			}

		} else {
			//package level raw messages won't have a Test value
			switch rawTestStep.Action {
			case "output":
				packageResult.Output = append(packageResult.Output, rawTestStep.Output)
				packageMap[rawTestStep.Package] = packageResult
			case "pass":
				packageResult.Result = rawTestStep.Action
				packageResult.Elapsed = rawTestStep.Elapsed
				packageMap[rawTestStep.Package] = packageResult
			case "fail":
				packageResult.Result = rawTestStep.Action
				packageResult.Elapsed = rawTestStep.Elapsed
				packageMap[rawTestStep.Package] = packageResult
			default:
				panic(fmt.Sprintf("unexpected action: %s", rawTestStep.Action))
			}
		}
	}

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	//transfer each test result map in each package result to a test result slice
	for j, packageResult := range packageMap {
		for _, testResults := range packageResult.TestMap {
			packageResult.Tests = append(packageResult.Tests, testResults...)
		}
		packageMap[j] = packageResult

		//clear test result map once all values transfered to slice - needed for testing so will check against an empty map
		for k := range packageMap[j].TestMap {
			delete(packageMap[j].TestMap, k)
		}
	}

	//sort all the test results in each package result slice - needed for testing so it's easy to compare ordered tests
	for _, pr := range packageMap {
		sort.SliceStable(pr.Tests, func(i, j int) bool {
			return pr.Tests[i].Test < pr.Tests[j].Test
		})
	}

	var testRun TestRun
	testRun.CommitDate = "Tue Sep 21 18:06:25 2021 -0700"
	testRun.CommitSha = "46baf6c6be29af9c040bc14195e195848598bbae"
	testRun.JobRunDate = "Tue Sep 21 21:06:25 2021 -0700"

	//add all the package results to the test run
	for _, pr := range packageMap {
		testRun.PackageResults = append(testRun.PackageResults, pr)
	}

	//sort all package results in the test run
	sort.SliceStable(testRun.PackageResults, func(i, j int) bool {
		return testRun.PackageResults[i].Package < testRun.PackageResults[j].Package
	})

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
		Results:    make(map[string]TestResult),
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
