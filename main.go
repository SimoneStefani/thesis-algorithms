package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/SimoneStefani/thesis-algorithms/structures/asl"
	"github.com/SimoneStefani/thesis-algorithms/structures/fastmt"
	"github.com/SimoneStefani/thesis-algorithms/structures/hashlist"
	"github.com/SimoneStefani/thesis-algorithms/structures/mt"
	. "github.com/SimoneStefani/thesis-algorithms/utilities"
)

func main() {

	// test := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	// sl, _ := asl.NewSkipList(test)
	// asl.PrintListAuthenticators(*sl)

	// //Print Examples for Searching the Skip List

	// // pos, exists := asl.Lookup(*sl, "k")
	// // if exists {
	// // 	fmt.Printf("%s is part of the list at pos: %d\n", "10", pos)
	// // } else {
	// // 	fmt.Printf("%t\n", exists)
	// // }

	// // for _, el := range test {
	// // 	pos, exists = asl.Lookup(*sl, el)
	// // 	if exists {
	// // 		fmt.Printf("%s is part of the list at pos: %d\n", el, pos)
	// // 	} else {
	// // 		fmt.Printf("%t\n", exists)
	// // 	}
	// // }
	// return

	// get absolute path of current folder
	basePath := GetPath()

	// parse the command line arguments
	algo, op, fileName, iter := ParseCommand()
	fmt.Printf("Running experiment with algo=%s and op=%s from %s...\n\n", *algo, *op, *fileName)

	// load data from specific file
	sourcePath := basePath + "/source/" + *fileName
	data := LoadData(sourcePath)

	// run experiment
	buildTimeResults, buildMemResults, veriTimeResults, veriMemResults := runExperiment(data, algo, *iter)

	// write to file the stringified result.
	// output file name pattern: result_[algo]_[inputName]
	// e.g. result_mt_uniform_samples_100.txt
	result := formatResults(buildTimeResults, buildMemResults, veriTimeResults, veriMemResults)
	resultName := "result_" + *algo + "_" + *fileName

	WriteData(basePath+"/results/"+resultName, result)
}

func formatResults(build_t []int64, build_m []int64, veri_t []int64, veri_m []int64) string {
	results := ""
	for i := 0; i < len(build_t); i++ {
		if i+1 == len(build_t) {
			results = results + strconv.FormatInt(build_t[i], 10) + ", " + strconv.FormatInt(build_m[i], 10) + ", " + strconv.FormatInt(veri_t[i], 10) + ", " + strconv.FormatInt(veri_m[i], 10)
		} else {
			results = results + strconv.FormatInt(build_t[i], 10) + ", " + strconv.FormatInt(build_m[i], 10) + ", " + strconv.FormatInt(veri_t[i], 10) + ", " + strconv.FormatInt(veri_m[i], 10) + "\n"
		}
	}
	return results
}

func runExperiment(data []string, algo *string, iter int) ([]int64, []int64, []int64, []int64) {

	buildTime, buildMem := runBuildExperiment(data, algo, iter)
	verificationTime, verificationMem := runVerificationExperiment(data, algo, iter)

	return buildTime, buildMem, verificationTime, verificationMem
}

func runBuildExperiment(data []string, algo *string, iter int) ([]int64, []int64) {
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
		} else if *algo == "sl" {
			start = time.Now()
			asl.NewSkipList(data)
			t = time.Now()
		}

		a := GetMemUsage()

		timeTrials = append(timeTrials, t.Sub(start).Nanoseconds())
		memTrials = append(memTrials, a-b)
	}
	return timeTrials, memTrials
}

func runVerificationExperiment(data []string, algo *string, iter int) ([]int64, []int64) {
	var timeTrials []int64
	var memTrials []int64
	averageTimePosition := len(data) / 2

	for i := 0; i < iter; i++ {
		var start time.Time
		var t time.Time
		var b int64

		runtime.GC()

		if *algo == "mt" {
			start = time.Now()
			root, path, _, _ := mt.VerifyTransaction(data[averageTimePosition], data)
			t = time.Now()

			runtime.GC()
			debug.SetGCPercent(-1)
			b = GetMemUsage()
			mt.CheckPath(data[averageTimePosition], root, path)
		} else if *algo == "hl" {
			start = time.Now()
			root, path, _, _ := hashlist.VerifyTransaction(data[averageTimePosition], data)
			t = time.Now()

			runtime.GC()
			debug.SetGCPercent(-1)
			b = GetMemUsage()
			hashlist.CheckPath(data[averageTimePosition], root, path)
		} else if *algo == "fmt" {
			start = time.Now()
			root, path, _, _ := fastmt.VerifyTransaction(data[averageTimePosition], data)
			t = time.Now()

			runtime.GC()
			debug.SetGCPercent(-1)
			b = GetMemUsage()
			fastmt.CheckPath(data[averageTimePosition], root, path)
		} else if *algo == "sl" {
			return nil, nil
		}

		a := GetMemUsage()

		timeTrials = append(timeTrials, t.Sub(start).Nanoseconds())
		memTrials = append(memTrials, a-b)
	}

	// fmt.Printf("%d transactions - Average time over %d samples: %v\n", len(data), iter, avg)

	return timeTrials, memTrials
}
