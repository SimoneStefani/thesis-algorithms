package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	data := loadData("./samples/uniform_samples_limited.txt")

	l := List{}
	l.BuildList(data)

	l.Show()
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
