# Obfuscatecsv
*Notes*
1. If mulitple columns have the same data, then they will be obfuscated to the same value to preserve identity of same value in two columns
2. The "prefix" is required and recommended to be something that is related to the data. For example, if names are being obfuscated, then use "name" as the prefix.
3. The sequences are simply the row and column of the first occurence of the value. That gives you a way to work backward if you need to.
4. The default delimiter between the row and column sequence number is a dash. If no delimiter is desired just use "" as shown below.

Use -help to show:
```
$ obfuscatecsv -help
Help Message

Usage: obfuscatecsv [options]
  -c string
    	Range spec for columns to obfuscate
  -d string
    	Delimiter for sequences (default "-")
  -headers
    	CSV has headers (default true)
  -help
    	Show help message
  -i string
    	Input CSV filename; default STDIN
  -keep
    	Keep CSV headers on output (default true)
  -o string
    	Output CSV filename; default STDOUT
  -prefix string
    	Prefix for obfuscator value
$
```

# Examples
Obfuscate first two columns:
```
$ cat test1.csv
A,B,C
abc,def,Army
def,abc,Navy
ijk,abc,Navy
zyz,def,Army
abc,abc,AF
$ obfuscatecsv -i test1.csv -prefix XT -c 1,2
A,B,C
XT2-0,XT2-1,Army
XT2-1,XT2-0,Navy
XT4-0,XT2-0,Navy
XT5-0,XT2-1,Army
XT2-0,XT2-0,AF
```
Chained/piped example that obfuscates all the columns, but with 
different prefixes.
```
$ obfuscatecsv -i test1.csv -prefix XT -c 1,2 -d "" | obfuscatecsv -prefix DOD -c 3
A,B,C
XT20,XT21,DOD2-2
XT21,XT20,DOD3-2
XT40,XT20,DOD3-2
XT50,XT21,DOD2-2
XT20,XT20,DOD6-2
$
```
