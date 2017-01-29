package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	cols := flag.String("c", "", "Order of columns from input")
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	help := flag.Bool("help", false, "Show usage message")
	flag.Parse()

	if *help {
		usage("Help Message")
	}

	if *cols == "" {
		usage("Missing new order of columns")
	}

	tokens := strings.Split(*cols, ",")
	outn := make([]int, len(tokens))

	for n := range tokens {
		i, err := strconv.Atoi(tokens[n])
		if err != nil {
			log.Fatalf("Value not a number:%v\n", tokens[n])
		}
		if i < 1 {
			log.Fatalf("Columns start at one:%v\n", tokens[n])
		}
		outn[n] = i
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
	outs := make([]string, len(tokens))

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
		if row == 0 {
			if *headers && *keep {
			} else {
				row++
				continue
			}
		}
		for n, m := range outn {
			outs[n] = cells[m-1]
		}
		err := w.Write(outs)
		if err != nil {
			log.Fatalf("csv.Write:\n%v\n", err)
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
