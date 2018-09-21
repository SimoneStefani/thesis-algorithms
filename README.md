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
  * `hl` = Hash List
  * `mt` = Merkle Tree
  * `fmt` = Fast Merkle Tree
  * `sl` = Authenticated Append-only Skip List (AASL)

* `-op` = the operation to perform
  * `build` = building the data structure and verifying a transaction

* `-name` =  the name of to the data source file
  * example: uniform_samples_100.txt

* `-iter` =  number of iterations

Full example:

```bash
./thesis -gopath=./ -algo=mt -op=build -name=uniform_samples_100.txt -iter=10
```

The output is written in a file in the `results` folder (don't remove the `.gitkeep` file). The name of the output file has the following pattern:

```
result_[algo]_[inputName]

e.g. result_mt_uniform_samples_100.txt
```

The content of the output files is layed out in the following form (where `,` is the separator) constituting a list of trials results:

```
[build_time], [build_memory], [avg_verification_time], [avg_verification_memory]\n
```

