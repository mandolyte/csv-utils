# Editcsv
This utility will edit a CSV and either update inline or add update 
as a new column.

Use -help to show:
```
$ editcsv -help
Help Message

Usage: editcsv [options] input.csv output.csv
  -add
    	Add replace string as a new column; default, replace in-place
  -addHeader string
    	Header to use for added column (default "ADDED")
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

## Examples
Put an "x-" in front of any cell value beginning with the letter "a".
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
```
Replace matches with a constant value, in this case "--elided--".
```
$ go run editcsv.go -pattern "^.*y$" -replace "--elided--" < test1.csv
A,B,C
abc,def,--elided--
one,two,--elided--
go,abacus,Marine
Android,Ubuntu,Linux
```
Replace matches (cell values in column 2 only) that end in letter "o",
adding a new column named "final" for the updated column 2.
```
$ editcsv -pattern "^.*o$" -replace "--elided--" -c 2 -add=true -addHeader "final" < test1.csv
A,B,C,final
abc,def,Army,def
one,two,Navy,--elided--
go,abacus,Marine,abacus
Android,Ubuntu,Linux,Ubuntu
```
