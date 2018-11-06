# Diffcsv

Todo:
- Considering adding a "map" option where you can supply a JSON formatted map of value transformations. Possibly using regualar expressions as matching values.

Changes (2018-11-06):
- Renamed the alias parameters as "alias1" and "alias2" with defaults of "f1" and "f2"
- Added `trimSpace` and `ignoreCase` options. The ignore case option applies to the *key* values.

Changes (2018-10-31):
- Added aliasing option of input files; default is DF1 and DF2 as before
- Added option to add numbers to column headers to make it easier to 
reference columns with differences

Use the -help argument to show:
```
$ diffcsv -help

Usage: diffcsv [options]
  -alias1 string
        Alias for first input file; default F1 (default "F1")
  -alias2 string
        Alias for second input file; default F2 (default "F2")
  -colnums
        Add difference column numbers to headers
  -f1 string
        First CSV file name to compare
  -f2 string
        Second CSV file name to compare
  -help
        Show help message
  -ignoreCase
        Ignore case when comparing; default true (default true)
  -key int
        Key column in input CSVs (first is 1); must be unique
  -noeq
        Suppress matches, showing only differences
  -o string
        Output CSV file for differences
  -ondupFirst
        On duplicate key, keep first one
  -ondupLast
        On duplicate key, keep last  one
  -trimSpace
        Ignore leading and trailing spaces when comparing; default true (default true)

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
	NOTE! if duplicate keys exist, then there are options to keep
	the first or to keep the last one. Default is to error out.
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

## Normal Cases

Compare two identical files (using same file for both inputs):
```
$ cat input1.csv
A,B,C
X,1,1
Y,2,2
Z,3,3
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input1.csv -o test1.csv
2018/10/08 06:46:55 Start: Oct  8 06:46:55.040
2018/10/08 06:46:55 Processing input #1:input1.csv
2018/10/08 06:46:55 Number of rows in file input1.csv:3
2018/10/08 06:46:55 Processing input #2:input1.csv
2018/10/08 06:46:55 Number of rows in file input1.csv:3
2018/10/08 06:46:55 Number of combined unique keys:3
2018/10/08 06:46:55 End: Oct  8 06:46:55.041
2018/10/08 06:46:55 Elapsed time 842.333µs
2018/10/08 06:46:55 ------- Summary -------
2018/10/08 06:46:55 Equal Count: 3
2018/10/08 06:46:55 Key Diff Count: 0
2018/10/08 06:46:55 Unique to input #1: 0
2018/10/08 06:46:55 Unique to input #2: 0
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
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input3.csv -o test3.csv
... elided ...
$ cat test2.csv
STATUS,A,B,C
EQ,X,1,1
EQ,Y,2,2
"DF1=2,3",Z,3,3
"DF2=2,3",Z,9,9
```

Same as above, but show only differences; use aliases and column numbers:
```
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input3.csv \
	-o test3.csv -noeq \
	-df1 i1 -df2 i2 -colnums
2018/10/31 07:03:26 Start: Oct 31 07:03:26.298
2018/10/31 07:03:26 Processing input #1:input1.csv
2018/10/31 07:03:26 Number of rows in file input1.csv:3
2018/10/31 07:03:26 Processing input #2:input3.csv
2018/10/31 07:03:26 Number of rows in file input3.csv:3
2018/10/31 07:03:26 Number of combined unique keys:3
2018/10/31 07:03:26 End: Oct 31 07:03:26.300
2018/10/31 07:03:26 Elapsed time 1.9993ms
2018/10/31 07:03:26 ------- Summary -------
2018/10/31 07:03:26 Equal Count: 2
2018/10/31 07:03:26 Key Diff Count: 1
2018/10/31 07:03:26 Unique to input #1: 0
2018/10/31 07:03:26 Unique to input #2: 0
$ cat test3.csv
STATUS,1-A,2-B,3-C
"i1=2,3",Z,3,3
"i2=2,3",Z,9,9
$ 
```

Compare two files where keys are not the same:
```
$ cat input2.csv
A,B,C
X,1,1
Y,2,2
W,3,3
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input2.csv -o test3.csv
... elided ...
$ cat test3.csv
STATUS,A,B,C
IN=2,W,3,3
EQ,X,1,1
EQ,Y,2,2
IN=1,Z,3,3
```

Compare two files where trim space and ignore case are needed:
```
$ cat input7.csv
A,B,C
X,1,1
Y,2,3
W,3,3
$ cat input8.csv
A,B,C
x,1,1
 Y ,2,3
 w ,3,3
$ go run diffcsv.go -key 1 -f1 input7.csv -f2 input8.csv -alias1 f1 -alias2 f2 -trimSpace=true -ignoreCase=false -o test7.csv
2018/11/06 13:37:09 Start: Nov  6 13:37:09.884
2018/11/06 13:37:09 Processing input #1:input7.csv
2018/11/06 13:37:09 Number of rows in file input7.csv:3
2018/11/06 13:37:09 Processing input #2:input8.csv
2018/11/06 13:37:09 Number of rows in file input8.csv:3
2018/11/06 13:37:09 Number of combined unique keys:5
2018/11/06 13:37:09 End: Nov  6 13:37:09.886
2018/11/06 13:37:09 Elapsed time 1.9977ms
2018/11/06 13:37:09 ------- Summary -------
2018/11/06 13:37:09 Equal Count: 1
2018/11/06 13:37:09 Key Diff Count: 0
2018/11/06 13:37:09 Unique to input #1: 2
2018/11/06 13:37:09 Unique to input #2: 2
$ cat test7.csv
STATUS,A,B,C
IN=f1,W,3,3
IN=f1,X,1,1
EQ,Y,2,3
IN=f2,w,3,3
IN=f2,x,1,1
```

## Error Conditions

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

Compare two files where one has a non-unique key:
```
$ cat input6.csv
A,B,C,D
X,1,1,1
Y,2,2,2
Z,9,9,9
X,1,2,3
$ go run diffcsv.go -key 1 -f1 input1.csv -f2 input6.csv -o test6.csv
2018/10/05 07:15:00 Start: Oct  5 07:15:00.105
2018/10/05 07:15:00 Processing input #1:input1.csv
2018/10/05 07:15:00 Number of rows in file input1.csv:3
2018/10/05 07:15:00 Processing input #2:input6.csv
2018/10/05 07:15:00 Key value not unique: X on row 4
exit status 1
$
```
