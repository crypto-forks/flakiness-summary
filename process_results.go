package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Event struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output"`
	Elapsed float64   `json:"Elapsed"`
}

type JobOutput struct {
	CommitSha  string                  `json:"CommitSha"`
	CommitDate time.Time               `json:"CommitDate"`
	JobStarted time.Time               `json:"JobStarted"`
	Results    map[string][]TestResult `json:"Results"`
}

type TestResult struct {
	Result  string  `json:"Result"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

func main() {
	commitDate, err := time.Parse(time.RFC3339, os.Getenv("COMMIT_DATE"))
	if err != nil {
		panic(err)
	}

	jobStarted, err := time.Parse(time.RFC3339, os.Getenv("JOB_STARTED"))
	if err != nil {
		panic(err)
	}

	jobOutput := JobOutput{
		CommitSha:  os.Getenv("COMMIT_SHA"),
		CommitDate: commitDate,
		JobRunDate: jobStarted,
		Results:    make(map[string][]TestResult),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			panic(err)
		}

		if event.Test == nil {
			continue
		}

		switch event.Action {
		case "run":
			// TODO
		case "pass":
			// TODO
		case "fail":
			// TODO
		case "output":
			// TODO
		default:
			panic(fmt.Sprintf("unexpected action: %s", event.Action))
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// TODO: write results to DB

	// TODO: output a flakiness summary
	// see https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-output-parameter
}
