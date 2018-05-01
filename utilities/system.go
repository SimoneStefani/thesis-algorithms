package utilities

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func ParseCommand() (*string, *string, *string, *int) {

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

func GetPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func GetMemUsage() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return int64(m.Alloc)
}
