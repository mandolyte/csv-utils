## Dedupcsv
This utility removed duplicate rows. The input must be sorted!

Use -help to show:
```
$ uniqcsv -help
Help Message

Usage: uniqcsv [options]
NOTE: must be sorted; only compares row against prior row.  -headers
        CSV has headers (default true)
  -help
        Show help message
  -i string
        Input CSV filename; default STDIN
  -keep
        Keep CSV headers on output (default true)
  -o string
        Output CSV filename; default STDOUT
```

For example:
```
$ cat test1.csv 
A,B,C
1,2,3
1,2,3
4,5,6
4,5,6
d,e,f
d,e,f
d,e,f
$ go run dedupcsv.go < test1.csv 
A,B,C
1,2,3
4,5,6
d,e,f
$
```
