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

var rs *rangespec.RangeSpec
var cs *rangespec.RangeSpec

func main() {
	rows := flag.String("r", "1-", "Range spec for rows")
	cols := flag.String("c", "1-", "Range spec for columns")
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	help := flag.Bool("help", false, "Show usage message")
	flag.Parse()

	if *help {
		usage("Help Message")
		os.Exit(0)
	}

	/* check parameters */
	if *rows == "" {
		usage("Required: Missing range specification for rows")
		os.Exit(0)
	}

	rs, rserr := rangespec.New(*rows)
	if rserr != nil {
		log.Fatalf("Invalid row range spec:%v, Error:\n%v\n", *rows, rserr)
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
			err := writeRow(w, cells, cs)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
			continue
		}
		row++
		if rs.InRange(row - 1) {
			err := writeRow(w, cells, cs)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
		}
		if row > rs.Max {
			break
		}
	}
	w.Flush()
}

func writeRow(w *csv.Writer, cells []string, cs *rangespec.RangeSpec) error {
	if cs == nil {
		err := w.Write(cells)
		if err != nil {
			return err
		}
		return nil
	}
	var outcells []string
	for m, c := range cells {
		if cs.InRange(uint64(m + 1)) {
			outcells = append(outcells, c)
		}
	}
	if len(outcells) == 0 {
		return fmt.Errorf("Column range outside actual columns:%v\n\n", cs)
	}
	err := w.Write(outcells)
	if err != nil {
		return err
	}
	return nil
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: splitcsv [options] input.csv output.csv\n")
	flag.PrintDefaults()
}
