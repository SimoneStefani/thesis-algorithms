package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	results := ""

	for i := 100; i < 1000; i += 50 {
		path := "./samples/uniform/uniform_samples_" + strconv.Itoa(i) + ".txt"
		data := loadData(path)
		results = results + strconv.FormatInt(evaluteOperations(data, 10), 10) + ", "
	}

	fmt.Println(results)
	writeData(results)
}

func evaluteOperations(data []string, iter int) int64 {
	var trials []int64

	for i := 0; i < iter; i++ {
		l := List{}
		start := time.Now()
		l.BuildList(data)
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

func writeData(data string) {
	path := "./out/results.txt"

	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	var file, e = os.OpenFile(path, os.O_RDWR, 0644)
	if e != nil {
		log.Fatal(e)
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
