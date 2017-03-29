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
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	help := flag.Bool("help", false, "Show usage message")
	force := flag.Bool("f", false, "Force concatenation of different width CSV files")
	flag.Parse()

	if *help {
		usage("Help Message")
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		usage("No files specified to concatenate!")
		os.Exit(0)
	}

	if !*headers {
		*keep = false
		log.Println("If no headers, keep option is auto-set to false")
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
	log.Println("Individual file row counts include header row")
	log.Println("Total row count does not include header rows")

	var total uint64
	var firstfilecolumncount int
	for n, f := range flag.Args() {
		// open input file
		var r *csv.Reader
		fi, fierr := os.Open(f)
		if fierr != nil {
			log.Fatal("os.Open() Error:" + fierr.Error())
		}
		defer fi.Close()
		r = csv.NewReader(fi)
		if n == 0 {
			r.FieldsPerRecord = 0
		} else {
			if *force {
				r.FieldsPerRecord = -1
			} else {
				r.FieldsPerRecord = firstfilecolumncount
			}
		}

		// read loop for CSV files
		var row uint64
		row = 0
		for {
			// read the csv file
			cells, rerr := r.Read()
			if rerr == io.EOF {
				break
			}
			if rerr != nil {
				log.Fatalf("csv.Read:\n%v\n", rerr)
			}
			if n == 0 && row == 0 {
				firstfilecolumncount = len(cells)
				if *headers {
					if *keep {
						err := w.Write(cells)
						if err != nil {
							log.Fatalf("csv.Write:\n%v\n", err)
						}
					}
				} else {
					err := w.Write(cells)
					if err != nil {
						log.Fatalf("csv.Write:\n%v\n", err)
					}
				}
				row++
				continue
			}
			if n > 0 && row == 0 {
				if *headers {
					row++
					continue // omit headers on all but first file
				}
			}
			row++
			err := w.Write(cells)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
		}
		log.Printf("File %v had %v rows", f, row)
		total += row
	}
	w.Flush()
	if *headers {
		total -= uint64(len(flag.Args())) // don't include header rows in counts
	}
	log.Printf("Total rows in output %v has %v rows", *output, total)
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: catcsv [options] input1.csv input2.csv ...\n")
	flag.PrintDefaults()
}
