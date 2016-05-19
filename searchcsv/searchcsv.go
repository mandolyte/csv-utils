package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.build.ge.com/210019484/rangespec"
)

var cs *rangespec.RangeSpec
var re *regexp.Regexp

func main() {
	pattern := flag.String("pattern", "", "Search pattern")
	cols := flag.String("c", "", "Range spec for columns")
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	regex := flag.Bool("re", false, "Search pattern is a regular expression")
	flag.Parse()

	/* check parameters */
	if *pattern == "" {
		usage("Required: Missing pattern for search")
		os.Exit(0)
	}

	if *regex {
		re = regexp.MustCompile(*pattern)
	}

	if *cols != "" {
		var cserr error
		cs, cserr = rangespec.New(*cols)
		if cserr != nil {
			log.Fatalf("Invalid column range spec:%v, Error:\n%v\n", *cols, cserr)
		}
	}

	if *keep {
		if !*headers {
			log.Fatal("Cannot keep headers you don't have!")
		}
	}
	// open output file
	var w *csv.Writer
	if *output == "" {
		w = csv.NewWriter(os.Stdout)
	} else {
		fo, foerr := os.Create(*output)
		if foerr != nil {
			log.Fatal("os.Create() Error:" + foerr.Error())
		}
		defer fo.Close()
		w = csv.NewWriter(fo)
	}

	// open input file
	var r *csv.Reader
	if *input == "" {
		r = csv.NewReader(os.Stdin)
	} else {
		fi, fierr := os.Open(*input)
		if fierr != nil {
			log.Fatal("os.Open() Error:" + fierr.Error())
		}
		defer fi.Close()
		r = csv.NewReader(fi)
	}

	// ignore expectations of fields per row
	r.FieldsPerRecord = -1

	// read loop for CSV
	var row uint64
	for {
		// read the csv file
		cells, rerr := r.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Fatalf("csv.Read:\n%v\n", rerr)
		}
		if (row == 0) && *headers && *keep {
			row = 1
			err := w.Write(cells)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
			continue
		}
		row++
		// test row/columns for a match
		if patternMatches(cells, *pattern) {
			err := w.Write(cells)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
		}
	}
	w.Flush()
}

func patternMatches(c []string, pattern string) bool {
	found := false
	for n, v := range c {
		if cs == nil {
			if re == nil {
				found = strings.Contains(v, pattern)
			} else {
				found = re.MatchString(v)
			}
		} else {
			if cs.InRange(uint64(n + 1)) {
				if re == nil {
					found = strings.Contains(v, pattern)
				} else {
					found = re.MatchString(v)
				}
			}
		}
		if found {
			return true
		}
	}
	return false
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: splitcsv [options] input.csv output.csv\n")
	flag.PrintDefaults()
}
