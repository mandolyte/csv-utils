# Catcsv
This utility will concatenate CSV files.

Use -help to show:
```
$ catcsv -help
Help Message

Usage: catcsv [options] input1.csv input2.csv ...
  -f    Force concatenation of different width CSV files
  -headers
        CSV has headers (default true)
  -help
        Show usage message
  -keep
        Keep CSV headers on output (default true)
  -o string
        Output CSV filename; default STDOUT
```

## Examples
This first example shows an error due to different number of columns
in the input files.
```
$ go run catcsv.go test1.csv test2.csv 
2017/12/01 09:18:16 Individual file row counts include header row
2017/12/01 09:18:16 Total row count does not include header rows
2017/12/01 09:18:16 File test1.csv had 4 rows
2017/12/01 09:18:16 csv.Read:
line 1, column 0: wrong number of fields in line
exit status 1
$
```
This example shows use of the force option to concatenate anyway.
```
$ go run catcsv.go -f test1.csv test2.csv 
2017/12/01 09:18:28 Individual file row counts include header row
2017/12/01 09:18:28 Total row count does not include header rows
2017/12/01 09:18:28 File test1.csv had 4 rows
2017/12/01 09:18:28 File test2.csv had 4 rows
A,B
1,2
3,4
5,6
1,2,3
4,5,6
7,8,9
2017/12/01 09:18:28 Total rows in output  has 6 rows
$ 
```