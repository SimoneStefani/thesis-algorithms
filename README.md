#### Royal Institute of Technology KTH - Stockholm
# Blockchain Algorithms

_This repository contains code written during the spring semester 2018 by Dean Rauschenbusch and Simone Stefani_

## Run the experiment

The source files with the sample data are located in the folder `source` (don't remove the `.gitkeep` file). The program knows which file to load based on the command line arguments.

In order to run the experiment first ensure that there is valid data in the `source` folder and then compile the Go code:

```bash
cd thesis-algorithms

go build -o thesis *.go
```

Then run the experiment. The program expects the following arguments:
* `-algo` the algorithm to run
    + `hl` = Hash List
    + `mt` = Merkle Tree
    + (`fmt` = Fast Merkle Tree)
    + (`bf` = Bloom Filter)

* `-op` = the operation to perform
    + `build` = building the data structure
    + (`verify` = verifying a transaction)
* `-name` =  the name of to the data source file
    + example: uniform_samples_100.txt

Full example:

```bash
./thesis -gopath=./ -algo=mt -op=build -name=uniform_samples_100.txt
```

The output is written in a file in the `results` folder (don't remove the `.gitkeep` file). The name of the output file has the following pattern:

```
result_[algo]_[inputName]

e.g. result_mt_uniform_samples_100.txt


Sample Content:

1603576\n
1518621\n
1661780\n
1942438\n
1404330\n
1488375\n
1470545\n
1458146\n
1549256\n
1727792
```