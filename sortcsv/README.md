# Sortcsv
This utility will sort a CSV file. However, it is done in-memory
and has limits.

Use the -help argument to show:

```
$ go run sortcsv.go -help
  -a	Sort ascending (default true) (default true)
  -c int
    	Column to sort (default 0)
  -headers
    	CSV has headers (default true)
  -help
    	Show help message
  -i string
    	CSV file name to sort; default STDIN
  -m string
    	Map to use instead of column values
  -o string
    	CSV output file name; default STDOUT
Map parameter is a filename of a JSON string
The JSON string maps column values to values to be
used for sorting. For example the following will
sort all values of NP first:

{
"NP": 0,
"*": 1
}

The asterisk is used to provide a default.
A default rule is required. If the default
mapping is also an asterisk, then the original
value is used.
```
Example:
```
$ cat test1.csv
A,B,C,D,E,F,G,H,I
1,1,1,1,1,1,1,1,1
2,2,2,2,2,2,2,2,2
3,3,3,3,3,3,3,3,3
4,4,4,4,4,4,4,4,4
5,5,5,5,5,5,5,5,5
6,6,6,6,6,6,6,6,6
7,7,7,7,7,7,7,7,7
8,8,8,8,8,8,8,8,8
9,9,9,9,9,9,9,9,9
$ splitcsv -c 4-6 -r 4-6 < test1.csv | sortcsv -a=false -c 2
D,E,F
6,6,6
5,5,5
4,4,4
$
```


