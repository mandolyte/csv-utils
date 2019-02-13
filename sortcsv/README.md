# Sortcsv
This utility will sort a CSV file. However, it is done in-memory
and has limits.

<<<<<<< HEAD
*To Do* 
- Remove the mapping option... not that useful. Consider another
CSV utility to map values
- Allow specification of a comma delimited list of column numbers
to sort on
- To with above, add some method of specifying ascending or 
descending for each of the columns provided. Default to ascending
for all of them.
- To be consistent with other utilities, have columns start at 1,
not zero.

## Information
=======
>>>>>>> ca08bc7df061d86e49603212210647fd9764a7dc
Use the -help argument to show:

```
$ go run sortcsv.go -help
<<<<<<< HEAD
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
=======
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
$
>>>>>>> ca08bc7df061d86e49603212210647fd9764a7dc
```
Example:
```
$ cat test1.csv
<<<<<<< HEAD
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
=======
A,B,C
1,2,3
4,1,0
2,1,2
3,3,1
3,3,2
$ go run sortcsv.go -c 1,3 -s d,a < test1.csv
A,B,C
4,1,0
3,3,1
3,3,2
2,1,2
1,2,3
$ go run sortcsv.go -c 1,3 -s d,d < test1.csv
A,B,C
4,1,0
3,3,2
3,3,1
2,1,2
1,2,3
>>>>>>> ca08bc7df061d86e49603212210647fd9764a7dc
$
```


