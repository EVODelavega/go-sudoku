# go-sudoku
Simple program to solve sudoku puzzles

### Example usage:

```
go run sudoku.go -path example.input
```

Should output

```
+|===|===|===||===|===|===||===|===|===|+
||   |   |   ||   | 8 |   ||   |   |   ||
||   | 1 |   ||   | 5 |   ||   | 7 |   ||
||   | 3 | 7 || 1 |   | 2 || 6 | 8 |   ||
||---|---|---||---|---|---||---|---|---||
|| 1 |   | 6 ||   |   |   || 5 |   | 3 ||
||   |   |   || 5 |   | 9 ||   |   |   ||
||   |   |   || 2 |   | 1 ||   |   |   ||
||---|---|---||---|---|---||---|---|---||
|| 5 |   | 3 ||   | 7 |   || 1 |   | 4 ||
||   |   | 8 ||   |   |   || 7 |   |   ||
||   | 4 |   ||   |   |   ||   | 9 |   ||
+|===|===|===||===|===|===||===|===|===|+
+|===|===|===||===|===|===||===|===|===|+
|| 6 | 5 | 9 || 7 | 8 | 4 || 2 | 3 | 1 ||
|| 8 | 1 | 2 || 3 | 5 | 6 || 4 | 7 | 9 ||
|| 4 | 3 | 7 || 1 | 9 | 2 || 6 | 8 | 5 ||
||---|---|---||---|---|---||---|---|---||
|| 1 | 9 | 6 || 8 | 4 | 7 || 5 | 2 | 3 ||
|| 2 | 7 | 4 || 5 | 3 | 9 || 8 | 1 | 6 ||
|| 3 | 8 | 5 || 2 | 6 | 1 || 9 | 4 | 7 ||
||---|---|---||---|---|---||---|---|---||
|| 5 | 2 | 3 || 9 | 7 | 8 || 1 | 6 | 4 ||
|| 9 | 6 | 8 || 4 | 1 | 3 || 7 | 5 | 2 ||
|| 7 | 4 | 1 || 6 | 2 | 5 || 3 | 9 | 8 ||
+|===|===|===||===|===|===||===|===|===|+
```

The example.input file uses `0` values to mark empty fields. Spaces can be used to the same effect: any character other than 1-9 will be interpreted as an empty field.

### Other options

* `-path`: tell the solver where to find a puzzle input file
* `-raw`: Instead of an input file, this option allows you to pass the puzzle as an argument. The format is the same as the one used in the input file.
* `-time`: Work in progress - this is probably going to change to a `-debug` flag, for now it just tells you how long it took to initialize, load, and solve the puzzle
* `-h`: Get some help

#### Runall script

Just run `./runall.sh -h` to see what the options are, to run all example input files, and keep the compiled binary, run the script with bot `b` and `k` flags.

#### Example puzzles

The hard and medium example puzzles were taken from http://websudoku.com/ 
The easy puzzle was taken from a July issue of the London Metro free newspaper. (12/07/2016?)
