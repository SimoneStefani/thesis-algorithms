package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	data100 := loadData("./samples/uniform_samples_limited_100.txt")

	evaluteOperations(data100, 1000)
}

func evaluteOperations(data []string, iter int) {
	var trials []int64

	for i := 0; i < iter; i++ {
		l := List{}
		start := time.Now()
		l.BuildList(data)
		t := time.Now()
		fmt.Println(t.Sub(start).Nanoseconds())
		trials = append(trials, t.Sub(start).Nanoseconds())
	}

	var sum int64
	for _, trial := range trials {
		sum = int64(sum) + int64(trial)
	}

	avg := sum / int64(iter)

	fmt.Printf("Average time over %d samples: %v\n", iter, avg)
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
