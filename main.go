package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	// parse the command line arguments
	algo, op, fileName := parseCommand()
	fmt.Printf("Running experiment with algo=%s and op=%s from %s...\n\n", *algo, *op, *fileName)

	// load data from specific file
	path := "./source/" + *fileName
	data := loadData(path)

	// run experiment
	results := evaluteOperations(data, algo, 10)

	// write to file the stringified result.
	// output file name pattern: result_[algo]_[inputName]
	// e.g. result_mt_uniform_samples_100.txt
	parsedResult := strconv.FormatInt(results, 10)
	resultName := "result_" + *algo + "_" + *fileName
	writeData("./results/"+resultName, parsedResult)
}

func evaluteOperations(data []string, algo *string, iter int) int64 {
	var trials []int64

	for i := 0; i < iter; i++ {
		var start time.Time
		var t time.Time

		if *algo == "mt" {
			start = time.Now()
			NewTree(data)
			t = time.Now()
		} else if *algo == "hl" {
			l := List{}
			start = time.Now()
			l.BuildList(data)
			t = time.Now()
		}

		trials = append(trials, t.Sub(start).Nanoseconds())
	}

	var sum int64
	for _, trial := range trials {
		sum = int64(sum) + int64(trial)
	}

	avg := sum / int64(iter)

	fmt.Printf("%d transactions - Average time over %d samples: %v\n", len(data), iter, avg)

	return avg
}

func parseCommand() (*string, *string, *string) {

	// Parse algorithm:
	// hl -> hashlist
	// mt -> Merkle tree (default)
	// fmt -> fast Merkle tree
	// bf -> Bloom's filter
	algorithm := flag.String("algo", "mt", "the algorithm to use")

	// Parse operation:
	// build -> build the data structure (default)
	// verify -> verification of block
	operation := flag.String("op", "build", "the operation to perform")

	// Parse output file name
	fileName := flag.String("name", "pew", "the name of the input file")

	flag.Parse()

	return algorithm, operation, fileName
}

func loadData(path string) []string {
	file, err := os.Open(path)

	// file opening error logging
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// load strings from file into slice
	var data []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	// scanner parsing error logging
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}

func writeData(path string, data string) {

	// detect if file exists
	var _, e = os.Stat(path)

	// remove file if it exists
	if os.IsExist(e) {
		var e = os.Remove(path)
		if e != nil {
			log.Fatal(e)
		}
	}

	// create file
	var file, err = os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// write the results
	_, err = file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}
