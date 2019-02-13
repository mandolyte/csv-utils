# Sortcsv
This utility will sort a CSV file. However, it is done in-memory
and has limits.

## Information
Use the -help argument to show:

```
$ go run sortcsv.go -help
  -c string
        Comma delimited list of columns to sort (default "1")
  -headers
        CSV has headers (default true)
  -help
        Show help message
  -i string
        CSV file name to sort; default STDIN
  -o string
        CSV output file name; default STDOUT
  -s string
        Comma delimited list of letters 'a' or 'd', for ascending or descending (default is ascending)
```

Example:
```
$ cat test1.csv 
A,B,C
1,2,3
4,1,0
2,1,2
3,3,1
3,3,2
$ go run sortcsv.go -c 1,3 -s a,d -i test1.csv 
A,B,C
1,2,3
2,1,2
3,3,2
3,3,1
4,1,0
$ go run sortcsv.go -c 1,3 -s a,a -i test1.csv 
A,B,C
1,2,3
2,1,2
3,3,1
3,3,2
4,1,0
$ $
```


