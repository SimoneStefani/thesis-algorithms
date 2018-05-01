package utilities

import (
	"bufio"
	"log"
	"os"
)

func LoadData(path string) []string {
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

func WriteData(path string, data string) {

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
