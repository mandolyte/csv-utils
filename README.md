# Comparecsv
This utility compares two CSV files. It uses Merkle Tree
concepts (hashes of each row are the basis of the compare logic).

The method used enables compares of large CSV files.
The only memory used is for the maps of hash values.

Use -help to show:
```
$ go run comparecsv.go -help
  -f1 string
    	First CSV file name to compare
  -f2 string
    	Second CSV file name to compare
  -help
    	Show help message
NOTE 1: Headers on the CSV files are expected.
NOTE 2: Duplicates are omitted in all outputs.
$ 
```

For example:
```
$ go run comparecsv.go -f1 test2.csv -f2 test3.csv 
2017/12/04 10:46:34 Start at 2017-12-04 10:46:34.39331915 -0500 EST m=+0.000330969
2017/12/04 10:46:34 Number of rows in file 1:3
2017/12/04 10:46:34 Number of rows in file 2:3
2017/12/04 10:46:34 Number of rows in both files:2
2017/12/04 10:46:34 Number of rows ONLY in file 2:1
2017/12/04 10:46:34 Number of rows ONLY in file 1:1
2017/12/04 10:46:34 End at 2017-12-04 10:46:34.394256297 -0500 EST m=+0.001268059
2017/12/04 10:46:34 Elapsed time 937.263µs
$ 
```

Performance test using the wine reviews public data set at https://www.kaggle.com/zynicide/wine-reviews/data:
```
$ go run comparecsv.go -f1 ~/data/winemag-data-130k-v2.csv -f2 ~/data/test1.csv 
2017/12/04 10:50:09 Start at 2017-12-04 10:50:09.775178465 -0500 EST m=+0.000293320
2017/12/04 10:50:12 Number of rows in file 1:129971
2017/12/04 10:50:16 Number of rows in file 2:129969
2017/12/04 10:50:16 Number of rows in both files:129968
2017/12/04 10:50:16 Number of rows ONLY in file 2:1
2017/12/04 10:50:19 Number of rows ONLY in file 1:3
2017/12/04 10:50:19 End at 2017-12-04 10:50:19.216067801 -0500 EST m=+9.441182871
2017/12/04 10:50:19 Elapsed time 9.440889815s
$ 
```