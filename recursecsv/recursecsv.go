package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	parent := flag.Int("parent", 1, "Parent column; default 1")
	child := flag.Int("child", 2, "Child column; default 2")
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	help := flag.Bool("help", false, "Show usage message")
	flag.Parse()

	if *help {
		usage("Help Message")
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
	r.FieldsPerRecord = 2

	// read loop for CSV to load into memory
	var row uint64
	pcol := *parent - 1
	ccol := *child - 1
	parents := make(map[string][]string)
	for {
		// read the csv file
		cells, rerr := r.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Fatalf("csv.Read:\n%v\n", rerr)
		}
		if row == 0 {
			recurseHeaders[1] = cells[pcol]
			recurseHeaders[2] = cells[ccol]
			row++
			continue
		}
		_, ok := parents[cells[pcol]]
		if ok {
			parents[cells[pcol]] = append(parents[cells[pcol]],cells[ccol])
		} else {
			parents[cells[pcol]] = make([]string,0)
			parents[cells[pcol]] = append(parents[cells[pcol]],cells[ccol])
		}
		row++
	}
	w.Flush()
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	flag.PrintDefaults()
	os.Exit(0)
}

var recurseHeaders []string

func init() {
	recurseHeaders = append(recurseHeaders, "level", "", "", "Path", "Cycle")
}
