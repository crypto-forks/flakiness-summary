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
	Output  string    `json:"Output"`
	Elapsed float64   `json:"Elapsed"`
}

func main() {
	fmt.Println("COMMIT_SHA:", os.Getenv("COMMIT_SHA"))
	fmt.Println("COMMIT_TIMESTAMP:", os.Getenv("COMMIT_TIMESTAMP"))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			panic(err)
		}
		fmt.Println(event)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// TODO: write results to DB

	// TODO: output a flakiness summary
	// see https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-output-parameter
}
