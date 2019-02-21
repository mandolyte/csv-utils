# CSV Utilities

This repo has a collection of CSV utilities to manipulate
CSV files. Here is a brief description of each. Each utility 
is in its own folder and has its own README.
- catcsv: concatenate two CSV files
- comparecsv: compare two CSV files
- dedupcsv: remove duplicates in a CSV file
- diffcsv: shows differences between two CSV files
- editcsv: alter contents of a CSV; regexp replace supported
- obfuscatecsv: obscures content in a regular fashion
- pivotcsv: do a pivot table operation
- recursecsv: recursively process hierarchical data; supports 
the Oracle list of hierarchical functions
- reordercsv: alters order of columns of a CSV file
- searchcsv: outputs matching rows of a CSV file; regexp 
supported
- sortcsv: sorts a CSV file
- splitcsv: splits a CSV by columns and/or rows
- transformcsv: using a "text/template", will transform a CSV 
by applying the template for each row

Each utility has its own README with examples.

To install `go get github.com/mandolyte/csv-utils`. 

Afterwards you can use `go install` to compile the ones of
interest or just use `go run`.

To install all of them: `sh build_all.sh`.

To Do:
- document recursedata.go