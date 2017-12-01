# Searchcsv
Use the -help argument to show:

```
$ searchcsv -help
Help Message

Usage: searchcsv [options]
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
  -v	Omit rather than include matched rows
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


