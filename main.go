package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type env struct {
	PathSamples string
	PathResults string
}

func main() {

	//Load environment
	env := loadEnv()

	results := ""

	for i := 100; i < 3000; i += 50 {
		path := env.PathSamples + strconv.Itoa(i) + ".txt"
		data := loadData(path)
		if i+50 >= 3000 {
			results = results + strconv.FormatInt(evaluteOperations(data, 10), 10)
		} else {
			results = results + strconv.FormatInt(evaluteOperations(data, 10), 10) + ","
		}
	}

	fmt.Println(results)
	writeData(env.PathResults+"results_mt2.txt", results)
}

func evaluteOperations(data []string, iter int) int64 {
	var trials []int64

	for i := 0; i < iter; i++ {
		// l := List{}
		start := time.Now()
		NewTree(data)
		// l.BuildList(data)
		t := time.Now()
		// fmt.Println(t.Sub(start).Nanoseconds())
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

func loadEnv() env {
	envPath := "./env.json"
	file, err1 := ioutil.ReadFile(envPath)
	if err1 != nil {
		fmt.Printf("error while reading file &s\n", envPath)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	var env env

	err2 := json.Unmarshal(file, &env)
	if err2 != nil {
		fmt.Println("error:", err2)
		os.Exit(1)
	}
	return env
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
