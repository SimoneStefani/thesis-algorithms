# Blockchain Algorithms

## Run the experiment

The source files with the sample data are located in the folder `source` (don't remote the `.gitkeep` file). The program knows which file to load based on the command line arguments.

In order to run the experiment first ensure that there is valid data in the `source` folder and then compile the Go code:

```bash
cd thesis-algorithms

go build -o thesis *.go
```

Then run the experiment. The program expects as arguments the algorithm to use, the operation to perform and the name of the data source file. For example:

```bash
./thesis -algo=mt -op=build -name=uniform_samples_100.txt
```

The output is written in a file in the `results` folder (don't remote the `.gitkeep` file). The name of the output file has the following pattern:

```
result_[algo]_[inputName]

e.g. result_mt_uniform_samples_100.txt
```
