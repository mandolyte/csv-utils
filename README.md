# csv-utils
Split (rows and columns), sort, and search

## To do
Plan to make a "editcsv" which will apply a Regex Replace.

## Splitcsv
Use the -help argument to show:

```
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

## Searchcsv
Use the -help argument to show:

```
$ go run searchcsv.go -help
Help Message

Usage: splitcsv [options] input.csv output.csv
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

## Editcsv
Use -help to show:
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
Example:
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