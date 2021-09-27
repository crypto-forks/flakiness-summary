package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestReadJson_1Count_AllPass(t *testing.T) {

	fptr := "./test/data/test-result-1-count-pass-crypto.json"

	f, err := os.Open(fptr)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		fmt.Println(s.Text())
	}

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
}
