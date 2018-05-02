package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/Daynex/thesis-algorithms/structures/fastmt"
	"github.com/Daynex/thesis-algorithms/structures/hashlist"
	"github.com/Daynex/thesis-algorithms/structures/mt"
	. "github.com/Daynex/thesis-algorithms/utilities"
)

func main() {

	// get absolute path of current folder
	basePath := GetPath()

	// parse the command line arguments
	algo, op, fileName, iter := ParseCommand()
	fmt.Printf("Running experiment with algo=%s and op=%s from %s...\n\n", *algo, *op, *fileName)

	// load data from specific file
	sourcePath := basePath + "/source/" + *fileName
	data := LoadData(sourcePath)

	// run experiment
	timeResults, memResults := evaluteOperations(data, algo, *iter)

	// write to file the stringified result.
	// output file name pattern: result_[algo]_[inputName]
	// e.g. result_mt_uniform_samples_100.txt
	parsedTimeResult := parseIntArrayToList(timeResults)
	parsedMemResult := parseIntArrayToList(memResults)
	timeResultName := "result_time_" + *algo + "_" + *fileName
	memResultName := "result_mem_" + *algo + "_" + *fileName

	WriteData(basePath+"/results/"+timeResultName, parsedTimeResult)
	WriteData(basePath+"/results/"+memResultName, parsedMemResult)
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

func evaluteOperations(data []string, algo *string, iter int) ([]int64, []int64) {
	var timeTrials []int64
	var memTrials []int64

	for i := 0; i < iter; i++ {
		var start time.Time
		var t time.Time

		runtime.GC()
		debug.SetGCPercent(-1)
		b := GetMemUsage()

		if *algo == "mt" {
			start = time.Now()
			mt.NewTree(data)
			t = time.Now()
		} else if *algo == "hl" {
			start = time.Now()
			hashlist.NewHashList(data)
			t = time.Now()
		} else if *algo == "fmt" {
			start = time.Now()
			fastmt.NewFastMerkleTree(data)
			t = time.Now()
		}

		a := GetMemUsage()

		timeTrials = append(timeTrials, t.Sub(start).Nanoseconds())
		memTrials = append(memTrials, a-b)
	}

	// fmt.Printf("%d transactions - Average time over %d samples: %v\n", len(data), iter, avg)

	return timeTrials, memTrials
}
