package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	//testResult.Results = append(testResult.Results, )
	//var results Result

	s := bufio.NewScanner(f)
	for s.Scan() {
		var testStep RawTestStep
		json.Unmarshal(s.Bytes(), &testStep)

		switch testStep.Action {
		case "run":
		case "pass":
		case "fail":
		case "output":
		default:
			panic(fmt.Sprintf("unexpected action: %s", testStep.Action))
		}
		//if testStep.Elapsed
		if testStep.Test == "" {
			if testStep.Action == "pass" || testStep.Action == "fail" {
				fmt.Print("Test: ")
				fmt.Print(testStep.Test)

				fmt.Print(" Action: ")
				fmt.Print(testStep.Action)

				fmt.Print(" Package: ")
				fmt.Print(testStep.Package)

				fmt.Print(" Elapsed: ")
				fmt.Println(testStep.Elapsed)
			}
		}
	}

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	var testResult1 TestResult
	testResult1.Elapsed = 0
	testResult1.Test = "TestSanitySha3_256"
	testResult1.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult1.Result = "pass"
	testResult1.Output = []string{"=== RUN   TestSanitySha3_256\n", "--- PASS: TestSanitySha3_256 (0.00s)\n"}

	var testResult2 TestResult
	testResult2.Elapsed = 0
	testResult2.Test = "TestSanitySha2_256"
	testResult2.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult2.Result = "pass"
	testResult2.Output = []string{"=== RUN   TestSanitySha2_256\n", "--- PASS: TestSanitySha2_256 (0.00s)\n"}

	var testResult3 TestResult
	testResult3.Elapsed = 0
	testResult3.Test = "TestSanitySha3_384"
	testResult3.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult3.Result = "pass"
	testResult3.Output = []string{"=== RUN   TestSanitySha3_384\n", "--- PASS: TestSanitySha3_384 (0.00s)\n"}

	var testResult4 TestResult
	testResult4.Elapsed = 0
	testResult4.Test = "TestSanitySha2_384"
	testResult4.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult4.Result = "pass"
	testResult4.Output = []string{"=== RUN   TestSanitySha2_384\n", "--- PASS: TestSanitySha2_384 (0.00s)\n"}

	var testResult5 TestResult
	testResult5.Elapsed = 0
	testResult5.Test = "TestSanityKmac128"
	testResult5.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult5.Result = "pass"
	testResult5.Output = []string{"=== RUN   TestSanityKmac128\n", "--- PASS: TestSanityKmac128 (0.00s)\n"}

	var testResult6 TestResult
	testResult6.Elapsed = 0
	testResult6.Test = "TestHashersAPI"
	testResult6.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult6.Result = "pass"
	testResult6.Output = []string{"=== RUN   TestHashersAPI\n", "    hash_test.go:114: math rand seed is 1632497249121800000\n", "--- PASS: TestHashersAPI (0.00s)\n"}

	var testResult7 TestResult
	testResult7.Elapsed = 0.23
	testResult7.Test = "TestSha3"
	testResult7.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult7.Result = "pass"
	testResult7.Output = []string{"=== RUN   TestSha3\n", "    hash_test.go:158: math rand seed is 1632497249122032000\n", "--- PASS: TestSha3 (0.23s)\n"}

	var testResult8 TestResult
	testResult8.Elapsed = 0.1
	testResult8.Test = "TestSha3/SHA3_256"
	testResult8.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult8.Result = "pass"
	testResult8.Output = []string{"=== RUN   TestSha3/SHA3_256\n", "    --- PASS: TestSha3/SHA3_256 (0.10s)\n"}

	var testResult9 TestResult
	testResult9.Elapsed = 0.12
	testResult9.Test = "TestSha3/SHA3_384"
	testResult9.Package = "github.com/onflow/flow-go/crypto/hash"
	testResult9.Result = "pass"
	testResult9.Output = []string{"=== RUN   TestSha3/SHA3_384\n", "    --- PASS: TestSha3/SHA3_384 (0.12s)\n"}

	var packageResult1 PackageResult
	packageResult1.Elapsed = 0.349
	packageResult1.Output = []string{"PASS\n", "ok  \tgithub.com/onflow/flow-go/crypto/hash\t0.349s\n"}
	packageResult1.Package = "github.com/onflow/flow-go/crypto/hash"
	packageResult1.Result = "pass"
	packageResult1.Tests = []TestResult{testResult1, testResult2, testResult3, testResult4, testResult5, testResult6, testResult7, testResult8, testResult9}

	var testRun TestRun
	//testResult.Results.
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
