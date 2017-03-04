package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mandolyte/rangespec"
)

var cs *rangespec.RangeSpec

func main() {
	prefix := flag.String("prefix", "", "Prefix for obfuscator value")
	cols := flag.String("c", "", "Range spec for columns to obfuscate")
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	delimiter := flag.String("d", "-", "Delimiter for sequences")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		usage("Help Message")
		os.Exit(0)
	}

	/* check parameters */
	if *prefix == "" {
		usage("Required: Missing prefix for obfuscation value")
		os.Exit(0)
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

	// Create value map to store mapping between
	// original values and obfuscated values
	valmap := make(map[string]string)

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
		//process(cells, *prefix, valmap, row, *width)
		for n, v := range cells {
			if cs.InRange(uint64(n + 1)) {
				obsv, ok := valmap[v]
				if ok {
					cells[n] = obsv
				} else {
					valmap[v] = fmt.Sprintf("%s%d%s%d", *prefix, row, *delimiter, n)
					cells[n] = valmap[v]
				}
			}
		}
		err := w.Write(cells)
		if err != nil {
			log.Fatalf("csv.Write:\n%v\n", err)
		}
	}
	w.Flush()
}

/*
func process(c []string, pf string, vm map[string]string, r uint64, w int) {

}
*/
func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: obfuscatecsv [options]\n")
	flag.PrintDefaults()
}
