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
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		usage("Help Message")
		os.Exit(0)
	}

	if !*headers {
		if *keep {
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
	var priorRow []string
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
			priorRow = make([]string,len(cells))
			_ = copy(priorRow,cells)
		}
		if (row == 0) && *headers && *keep {
			row = 1
			err := w.Write(cells)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
			continue
		}
		
		areEqual := testEq(priorRow,cells)
		if areEqual {
			continue
		}
		
		err := w.Write(cells)
		if err != nil {
			log.Fatalf("csv.Write:\n%v\n", err)
		}
		
		priorRow = make([]string,len(cells))
		_ = copy(priorRow,cells)

		row++
	}
	w.Flush()
}

func testEq(a, b []string) bool {

    if a == nil && b == nil { 
        return true; 
    }

    if a == nil || b == nil { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}


func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: uniqcsv [options]\n")
	fmt.Print("NOTE: must be sorted; only compares row against prior row.")
	flag.PrintDefaults()
}
