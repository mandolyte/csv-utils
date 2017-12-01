# Pivotcsv
Use -help to show:
```
$ pivotcsv -help
  -c int
    	Column to pivot (REQUIRED)
  -headers
    	CSV must have headers; cannot be false (default true)
  -help
    	Show help message
  -i string
    	CSV file name to pivot; default STDIN
  -nf string
    	Format to use for numbers (default "%v")
  -nv string
    	String to signal novalue; default is empty string
  -o string
    	CSV output file name; default STDOUT
  -on
    	Only consider numeric data and sum them (default true)
  -os
    	Consider data as strings and concatenate
  -s int
    	Column to sum/concat (REQUIRED)
  -sd string
    	Concatenation delimiter; default is comma (default ",")
```
## Examples
```
$ cat test1.csv
A,B,C,D,E,F
a1,b1,c1,d1,X,1
a2,b2,c2,d2,Y,3
a1,b1,c1,d1,X,3
a2,b2,c2,d2,Y,3
$ go run pivotcsv.go -i test1.csv -c 5 -s 6
A,B,C,D,X,Y
a1,b1,c1,d1,4,
a2,b2,c2,d2,,6

$ go run pivotcsv.go -i test1.csv -c 1 -s 2 -os
a1,a2,C,D,E,F
b1,,X,1
b1,,X,3
,"b2,b2",Y,3
$ cat test1.csv
A,B,C,D,E,F
a1,b1,c1,d1,X,1
a2,b2,c2,d2,Y,3
a1,b1,c1,d1,X,3
a2,b2,c2,d2,Y,3

$ cat test2.csv
A,B,C,D,E,F
a1,b1,c1,d1,X,1
a2,b2,c2,d2,X,3
a1,b1,c1,d1,X,3
a2,b2,c2,d2,X,3
a1,b1,c1,d1,Y,2
a2,b2,c2,d2,Y,4
a1,b1,c1,d1,Y,4
a2,b2,c2,d2,Y,4
a1,b1,c1,d1,Z,3
a2,b2,c2,d2,Z,5
a1,b1,c1,d1,Z,5
a2,b2,c2,d2,Z,5
$ go run pivotcsv.go -c 5 -s 6 < test2.csv
A,B,C,D,X,Y,Z
a1,b1,c1,d1,4,6,8
a2,b2,c2,d2,6,8,10
$

```
