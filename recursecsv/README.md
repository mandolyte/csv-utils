# Recursecsv
*Notes*
1. It will always output the normal hierarchical columns in this order: level, root, parent, child, path, and cycle (a Yes/No)
2. Note defaults shown in the help message below
3. At present it can only take two columns of data, the parent and child columns. If these have other associated values, they will have to be added back in to this output.
4. The input must have column headers, since they are re-used in the output CSV.


Use -help to show:
```
$ recursecsv -help
Help Message

  -child int
    	Child column; default 2 (default 2)
  -delimiter string
    	String for path delimiter (default ">")
  -help
    	Show usage message
  -i string
    	Input CSV filename; default STDIN
  -o string
    	Output CSV filename; default STDOUT
  -parent int
    	Parent column; default 1 (default 1)
  -start string
    	Start value of hierarchy
```

## Examples
Example with a cyclic condition.
```
$ cat test1.csv
parent,child
A,X
A,B
B,C
D,E
C,D
X,Y
Y,B
E,C
$ recursecsv -i test1.csv -start A
2017/12/01 09:56:39 Start at 2017-12-01 14:56:39.064464694 +0000 UTC
2017/12/01 09:56:39 Data loaded and ready to start recursing
2017/12/01 09:56:39 Working on A
2017/12/01 09:56:39 . elasped 66.33µs
2017/12/01 09:56:39 End at 2017-12-01 14:56:39.087153217 +0000 UTC
2017/12/01 09:56:39 Elapsed time 22.688732ms
Level,Root,parent,child,Path,Leaf,Cycle
1,A,A,B,>A>B>,No,No
2,A,B,C,>A>B>C>,No,No
3,A,C,D,>A>B>C>D>,No,No
4,A,D,E,>A>B>C>D>E>,No,No
5,A,E,C,>A>B>C>D>E>C>,No,Yes
1,A,A,X,>A>X>,No,No
2,A,X,Y,>A>X>Y>,No,No
3,A,Y,B,>A>X>Y>B>,No,No
4,A,B,C,>A>X>Y>B>C>,No,No
5,A,C,D,>A>X>Y>B>C>D>,No,No
6,A,D,E,>A>X>Y>B>C>D>E>,No,No
7,A,E,C,>A>X>Y>B>C>D>E>C>,No,Yes
$ 
```
Simple no cycle test.
```
$ recursecsv -i test2.csv -start A
2017/12/01 09:58:39 Start at 2017-12-01 14:58:39.319162864 +0000 UTC
2017/12/01 09:58:39 Data loaded and ready to start recursing
2017/12/01 09:58:39 Working on A
2017/12/01 09:58:39 . elasped 87.756µs
2017/12/01 09:58:39 End at 2017-12-01 14:58:39.319813 +0000 UTC
2017/12/01 09:58:39 Elapsed time 650.482µs
Level,Root,parent,child,Path,Leaf,Cycle
1,A,A,B,>A>B>,No,No
2,A,B,C,>A>B>C>,No,No
3,A,C,D,>A>B>C>D>,No,No
4,A,D,E,>A>B>C>D>E>,Yes,No
```
