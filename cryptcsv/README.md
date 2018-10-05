# Cryptcsv

This routine will encrypt/decrypt the selected column(s) using the supplied key.

Use -help to show:
```
$ cryptcsv -help
Help Message

Usage: cryptcsv [options]
  -c string
        Range spec for columns to obfuscate
  -d string
        Decrpytion key; required if decrypting
  -e string
        Encrpytion key; required if encrypting
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
$ 
```

# Examples
Encrypt the first and last columns:
```
$ cat test1.csv
A,B,C
abc,def,Army
def,abc,Navy
ijk,abc,Navy
zyz,def,Army
abc,abc,AF
$ cryptcsv -i test1.csv -c 1,3 -e abcdef -o test1-encrypted.csv
$ cat test1-encrypted.csv
$ cat test1-encrypted.csv 
A,B,C
LjZwW4XHoiXeg/5S9PItOmw7LQ==,def,cIJhrzIIYEXgAbcTPKEYpsKIJcw=
FaKIeKSORfKhZO+Sm3Rg3vEQKQ==,abc,/1WXOfya+LjAHWB2xr4zqo8Qmks=
EmNIdOIqir9TiT4mAf6o1vFYrQ==,abc,zao7y8CJgzW+G1ZSjRWelhIzNhw=
uowhS1km7U7B7k+aa8bWz0lUgw==,def,ShIYZBMV+PFG8JTud/FFRVGjtVQ=
V8WKIWunjW12OKC+MCcqlZqH2w==,abc,4ydx/qW9LierW6pQFeILRRtV
$
```

Now decrypt just the last column:
```
$ cryptcsv -i test1-encrypted.csv -o test1-decrypted.csv -c 3 -d abcdef
$ cat test1-decrypted.csv
,B,C
LjZwW4XHoiXeg/5S9PItOmw7LQ==,def,Army
FaKIeKSORfKhZO+Sm3Rg3vEQKQ==,abc,Navy
EmNIdOIqir9TiT4mAf6o1vFYrQ==,abc,Navy
uowhS1km7U7B7k+aa8bWz0lUgw==,def,Army
V8WKIWunjW12OKC+MCcqlZqH2w==,abc,AF
$
$ cksum test1.csv test1-decrypted.csv
```

Now decrypt both:
```
$ cryptcsv -i test1-encrypted.csv -o test1-decrypted-both.csv -d abcdef -c 1,3 
$ cat test1-decrypted-both.csv 
A,B,C
abc,def,Army
def,abc,Navy
ijk,abc,Navy
zyz,def,Army
abc,abc,AF
$ cksum test1.csv test1-decrypted-both.csv 
2235581246 69 test1.csv
2235581246 69 test1-decrypted-both.csv
$
``` 
