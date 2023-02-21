# Splitcsv
Use the -help argument to show:
```
$ go run splitcsv.go -help
Help Message

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
$ go run splitcsv.go -c 4-6 -r 4-6 < test1.csv
D,E,F
4,4,4
5,5,5
6,6,6
$
```

To upgrade to the new mod system:

1. created a subfolder named rangespec
2. copied the rangespec.go from project into it.
3. ran the "go mod" command:
```
$ go mod init github.com/mandolyte/csv-utils/splitcsv
```
4. then changed my import to be:
```go
import (
        "encoding/csv"
        "flag"
        "fmt"
        "io"
        "log"
        "os"
        "github.com/mandolyte/csv-utils/splitcsv/rangespec"
)
```

