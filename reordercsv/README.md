# Reordercsv
Use -help to show:
```
$ reordercsv -help
Help Message

  -c string
    	Order of columns from input
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
$
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
$ reordercsv -i test1.csv -c 3,2,1,1
C,B,A,A
1,1,1,1
2,2,2,2
3,3,3,3
4,4,4,4
5,5,5,5
6,6,6,6
7,7,7,7
8,8,8,8
9,9,9,9
$
```

