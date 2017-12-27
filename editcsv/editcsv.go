package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/mandolyte/csv-utils"
)

var cs *rangespec.RangeSpec
var re *regexp.Regexp

func main() {
	pattern := flag.String("pattern", "", "Search pattern")
	replace := flag.String("replace", "", "Regexp replace expression")
	addHdr := flag.String("addHeader", "ADDED", "Header to use for added column")
	cols := flag.String("c", "", "Range spec for columns")
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	help := flag.Bool("help", false, "Show help message")
	add := flag.Bool("add", false, "Add replace string as a new column; default, replace in-place")
	flag.Parse()

	if *help {
		usage("Help Message")
		os.Exit(0)
	}

	/* check parameters */
	if *replace == "" {
		usage("Required: Missing replace expression")
		os.Exit(0)
	}

	if *pattern == "" {
		usage("Required: Missing search expression")
		os.Exit(0)
	}
	re = regexp.MustCompile(*pattern)

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
			if *add {
				cells = append(cells, *addHdr)
			}
			err := w.Write(cells)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
			continue
		}
		row++
		// test row/columns for a match
		err := w.Write(patternMatches(cells, re, *replace, *add))
		if err != nil {
			log.Fatalf("csv.Write:\n%v\n", err)
		}
	}
	w.Flush()
}

func patternMatches(c []string, re *regexp.Regexp, replace string, add bool) []string {
	for n := range c {
		if cs == nil {
			newstring := re.ReplaceAllString(c[n], replace)
			if add {
				c = append(c, newstring)
			} else {
				c[n] = newstring
			}
		} else {
			if cs.InRange(uint64(n + 1)) {
				newstring := re.ReplaceAllString(c[n], replace)
				if add {
					c = append(c, newstring)
				} else {
					c[n] = newstring
				}
			}
		}
	}
	return c
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: editcsv [options] input.csv output.csv\n")
	flag.PrintDefaults()
}
