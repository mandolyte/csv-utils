# Transformcsv
This utility will take an input CSV and transform it using a text template.
The template is applied to every row in the CSV. The column headers are 
required. The column header names are used as map keys to the values
used by the template.

Use the -help argument to show:
```
$ go run transformcsv.go -help
Help Message

  -help
        Show usage message
  -i string
        Input CSV filename; default STDIN
  -m string
        Name of map in template; default is m (default "m")
  -o string
        Output filename; default STDOUT
  -t string
        Template to use for transformation
$ 
```

Given template:
```
$ cat template1.txt 
INSERT INTO atable (column1, column2, column3)
VALUES ('{{index .mp "column1"}}', '{{index .mp "column2"}}', '{{index .mp "column3"}}')
;
```

Given input CSV:
```
$ cat test1.csv 
column1,column2,column3
v1.1,v1.2,v1.3
v2.1,v2.1,v2.3
$ 
```

Then this command will generate SQL INSERT statements for each row
in the CSV file.
```
$ go run transformcsv.go -i test1.csv -t template1.txt -m mp -o trans1.sql
$ cat trans1.sql
INSERT INTO atable (column1, column2, column3)
VALUES ('v1.1', 'v1.2', 'v1.3')
;
INSERT INTO atable (column1, column2, column3)
VALUES ('v2.1', 'v2.1', 'v2.3')
;
$ 
```