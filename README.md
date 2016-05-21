# csv-utils
Split (rows and columns), sort, and search

## To do
Plan to make a "editcsv" which will apply a Regex Replace.

## Splitcsv
Use the -help argument to show:
```
$ go run splitcsv.go -help
Help Message

Usage: splitcsv [options] input.csv output.csv
  -c string
    	Range spec for columns
  -headers
    	CSV has headers (default true)
  -help
    	Show usage message
  -i string
    	Input CSV filename; default STDIN
  -keep
    	Keep CSV headers on output (default true)
  -o string
    	Output CSV filename; default STDOUT
  -r string
    	Range spec for rows
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
$ go run splitcsv.go -c 4-6 -r 4-6 < test1.csv
D,E,F
4,4,4
5,5,5
6,6,6
$
```

## Sortcsv
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
$ splitcsv -c 4-6 -r 4-6 < test1.csv | sortcsv -a=false -c 2
D,E,F
6,6,6
5,5,5
4,4,4
$
```

## Searchcsv
Use the -help argument to show:

```
$ go run searchcsv.go -help
Help Message

Usage: searchcsv [options] input.csv output.csv
  -c string
      Range spec for columns
  -headers
      CSV has headers (default true)
  -help
      Show help message
  -i string
      Input CSV filename; default STDIN
  -keep
      Keep CSV headers on output (default true)
  -o string
      Output CSV filename; default STDOUT
  -pattern string
      Search pattern
  -re
      Search pattern is a regular expression

```
Examples:
```
$ cat test1.csv
A,B,C
abc,def,Army
one,two,Navy
go,abacus,Marine
Android,Ubuntu,Linux
$ searchcsv -c 1 -pattern "y$" < test1.csv
A,B,C
$ searchcsv -c 3 -pattern "y$" < test1.csv
A,B,C
$ searchcsv -c 3 -pattern "y$" -re=true < test1.csv
A,B,C
abc,def,Army
one,two,Navy
$ searchcsv -c 3 -pattern "[mu][xy]$" -re=true < test1.csv
A,B,C
abc,def,Army
Android,Ubuntu,Linux
$ searchcsv -v -c 3 -pattern "[mu][xy]$" -re=true < test1.csv
A,B,C
one,two,Navy
go,abacus,Marine
```
## Editcsv
Use -help to show:
```
$ go run editcsv.go -help
Help Message

Usage: editcsv [options] input.csv output.csv
  -c string
      Range spec for columns
  -headers
      CSV has headers (default true)
  -help
      Show help message
  -i string
      Input CSV filename; default STDIN
  -keep
      Keep CSV headers on output (default true)
  -o string
      Output CSV filename; default STDOUT
  -pattern string
      Search pattern
  -replace string
      Regexp replace expression
```
Examples:
```
$ cat test1.csv
A,B,C
abc,def,Army
one,two,Navy
go,abacus,Marine
Android,Ubuntu,Linux
$ go run editcsv.go -pattern "^(a)" -replace "x-$1" < test1.csv
A,B,C
x-bc,def,Army
one,two,Navy
go,x-bacus,Marine
Android,Ubuntu,Linux
$ go run editcsv.go -pattern "^.*y$" -replace "--elided--" < test1.csv
A,B,C
abc,def,--elided--
one,two,--elided--
go,abacus,Marine
Android,Ubuntu,Linux
```
