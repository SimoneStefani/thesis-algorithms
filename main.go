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
	for i := 100; i < 3000; i += 50 {
		path := "./samples/uniform_samples_" + strconv.Itoa(i) + ".txt"
		data := loadData(path)
		evaluteOperations(data, 10)
	}
}

func evaluteOperations(data []string, iter int) {
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
