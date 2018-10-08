# Comparecsv
This utility compares two CSV files using Merkle Tree conceps, 
namely, hashes of the rows are the basis of the compare logic.

It is written to enable large CSV file comparisons. The only
memory consumed are maps of the row hashes.

Use -help to show:
```
$ go run comparecsv.go -help
  -f1 string
    	First CSV file name to compare
  -f2 string
    	Second CSV file name to compare
  -help
    	Show help message
NOTE 1: Headers on the CSV files are expected.
NOTE 2: Duplicates are omitted in all outputs.
$
```

It produces three output files, which are currently fixed:
- f1only.csv contains the rows unique to file 1
- f2only.csv contains the rows unique to file 2
- both.csv contains the rows common to both input files

## Examples
A simple test to validate basic operations:
```
$ go run comparecsv.go -f1 test2.csv -f2 test3.csv 
2017/12/04 11:15:29 Start at 2017-12-04 11:15:29.853501341 -0500 EST m=+0.000326007
2017/12/04 11:15:29 Number of rows in file 1:3
2017/12/04 11:15:29 Number of rows in file 2:3
2017/12/04 11:15:29 Number of rows in both files:2
2017/12/04 11:15:29 Number of rows ONLY in file 2:1
2017/12/04 11:15:29 Number of rows ONLY in file 1:1
2017/12/04 11:15:29 End at 2017-12-04 11:15:29.85432992 -0500 EST m=+0.001154546
2017/12/04 11:15:29 Elapsed time 828.715µs
$
```

A performance test using wine review public data set at
https://www.kaggle.com/zynicide/wine-reviews/data. Minor
changes are made to the original to make test1.csv.
```
$ comparecsv -f1 winemag-data-130k-v2.csv -f2 test1.csv 
2017/12/04 11:18:40 Start at 2017-12-04 11:18:40.631915938 -0500 EST m=+0.000781184
2017/12/04 11:18:43 Number of rows in file 1:129971
2017/12/04 11:18:49 Number of rows in file 2:129969
2017/12/04 11:18:49 Number of rows in both files:129968
2017/12/04 11:18:49 Number of rows ONLY in file 2:1
2017/12/04 11:18:51 Number of rows ONLY in file 1:3
2017/12/04 11:18:51 End at 2017-12-04 11:18:51.356633528 -0500 EST m=+10.725498483
2017/12/04 11:18:51 Elapsed time 10.72471747s
$ 
```