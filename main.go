package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {

	// get absolute path of current folder
	basePath := getPath()

	// parse the command line arguments
	algo, op, fileName, iter := parseCommand()
	fmt.Printf("Running experiment with algo=%s and op=%s from %s...\n\n", *algo, *op, *fileName)

	// load data from specific file
	sourcePath := basePath + "/source/" + *fileName
	data := loadData(sourcePath)

	// run experiment
	results := evaluteOperations(data, algo, *iter)

	// write to file the stringified result.
	// output file name pattern: result_[algo]_[inputName]
	// e.g. result_mt_uniform_samples_100.txt
	parsedResult := parseIntArrayToList(results)
	resultName := "result_" + *algo + "_" + *fileName
	destPath := basePath + "/results/" + resultName
	writeData(destPath, parsedResult)
}

func parseIntArrayToList(data []int64) string {
	results := ""
	for i := 0; i < len(data); i++ {
		if i+1 == len(data) {
			results = results + strconv.FormatInt(data[i], 10)
		} else {
			results = results + strconv.FormatInt(data[i], 10) + "\n"
		}
	}
	return results
}

func evaluteOperations(data []string, algo *string, iter int) []int64 {
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

	//fmt.Printf("%d transactions - Average time over %d samples: %v\n", len(data), iter, avg)

	return trials
}

func parseCommand() (*string, *string, *string, *int) {

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

	// Parse output file name
	iterations := flag.Int("iter", 10, "number of iterations")

	flag.Parse()

	return algorithm, operation, fileName, iterations
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

func getPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	return dir
}
