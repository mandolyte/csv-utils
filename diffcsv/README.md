# Diffcsv
Use the -help argument to show:
```
$ go run diffcsv.go -help

Usage: diffcsv [options]
  -f1 string
        First CSV file name to compare
  -f2 string
        Second CSV file name to compare
  -help
        Show help message
  -key int
        Key column in input CSVs (first is 1); must be unique
  -o string
        Output CSV file for differences

        Detailed Help:
        Inputs:
                - a key column
                - two input filenames
                - an output filename
        There will be two input files to compare and there will be
        one output file created:
        a) The first file will be read and stored into a map
        b) The second file will be read and stored into a map
        c) It is an error if a file has the same key value on two rows.
        Keys must be unique within each file.
        Note that key column number is one based, not zero based!
        d) Then all keys from both inputs are combined/deduped/sorted
        e) Then we range over the combined keyset and output a new CSV
        that has a new status column as the first column and the other columns
        from the inputs as the remaining columns.
        f) the new status column has the following values:
        - EQ meaning that the values for the key are same in both input files
        - IN=1 meaning that the key and values are only in input file #1
        - IN=2 similar for input file #2
        - DFn=x,y,..,z where n is either 1 or 2; followed by a comma delimited
        list of column numbers where the values for the key do not match.
        Note that the DF statuses always come in pairs, one for each input file.
        g) Limitations:
        - both input files must have the same number of columns
        - both must have a header row and the headers must be the same
```

Compare two identical files (using same file for both inputs):
```
$ cat input1.csv
A,B,C
X,1,1
Y,2,2
Z,3,3
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input1.csv -o test1.csv
2018/10/04 21:12:09 Start: Oct  4 21:12:09.833
2018/10/04 21:12:09 Number of rows in file input1.csv:3
2018/10/04 21:12:09 Number of rows in file input1.csv:3
2018/10/04 21:12:09 Number of combined unique keys:3
2018/10/04 21:12:09 End: Oct  4 21:12:09.835
2018/10/04 21:12:09 Elapsed time 1.9988ms
$ cat test1.csv
STATUS,A,B,C
EQ,X,1,1
EQ,Y,2,2
EQ,Z,3,3
```

Compare two files where keys are ok, but values are different:
```
$ cat input3.csv
A,B,C
X,1,1
Y,2,2
Z,9,9
$ cat test2.csv
STATUS,A,B,C
EQ,X,1,1
EQ,Y,2,2
"DF1=2,3",Z,3,3
"DF2=2,3",Z,9,9
```

Compare two files where keys are not the same:
```
$ cat input2.csv
A,B,C
X,1,1
Y,2,2
W,3,3
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input2.csv -o test3.csv
2018/10/04 21:22:24 Start: Oct  4 21:22:24.618
2018/10/04 21:22:24 Number of rows in file input1.csv:3
2018/10/04 21:22:24 Number of rows in file input2.csv:3
2018/10/04 21:22:24 Number of combined unique keys:4
2018/10/04 21:22:24 End: Oct  4 21:22:24.620
2018/10/04 21:22:24 Elapsed time 1.9997ms
$ cat test3.csv
STATUS,A,B,C
IN=2,W,3,3
EQ,X,1,1
EQ,Y,2,2
IN=1,Z,3,3
```

Compare two files with headers that don't match:
```
$ cat input4.csv
A,B,D
X,1,1
Y,2,2
Z,9,9
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input4.csv -o test4.csv
2018/10/04 21:25:36 Start: Oct  4 21:25:36.905
2018/10/04 21:25:36 Headers are not the same on input files
exit status 1
$
```

Compare two files that don't the same number of columns:
```
$ cat input5.csv
A,B,C,D
X,1,1,1
Y,2,2,2
Z,9,9,9
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input5.csv -o test5.csv
2018/10/04 21:27:24 Start: Oct  4 21:27:24.851
2018/10/04 21:27:24 Different number of columns:3 vs. 4
exit status 1
$
```