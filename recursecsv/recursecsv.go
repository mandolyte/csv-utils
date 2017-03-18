package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var w *csv.Writer

var parent = flag.Int("parent", 1, "Parent column; default 1")
var child = flag.Int("child", 2, "Child column; default 2")
var start = flag.String("start", "", "Start value of hierarchy")
var delimiter = flag.String("delimiter", ">", "String for path delimiter")
var input = flag.String("i", "", "Input CSV filename; default STDIN")
var output = flag.String("o", "", "Output CSV filename; default STDOUT")
var help = flag.Bool("help", false, "Show usage message")

func main() {
	flag.Parse()

	if *help {
		usage("Help Message")
	}

	// open output file
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
			log.Fatalf("csv.Read [row %v]:\n%v\n", row, rerr)
		}
		if row == 0 {
			recurseHeaders[1] = cells[pcol]
			recurseHeaders[2] = cells[ccol]
			writeRow(recurseHeaders[0], recurseHeaders[1],
				recurseHeaders[2], recurseHeaders[3], recurseHeaders[4])
			row++
			continue
		}
		_, ok := parents[cells[pcol]]
		if ok {
			parents[cells[pcol]] = append(parents[cells[pcol]], cells[ccol])
		} else {
			parents[cells[pcol]] = make([]string, 0)
			parents[cells[pcol]] = append(parents[cells[pcol]], cells[ccol])
		}
		row++
	}

	recurse(0, *start, *delimiter+*start, parents)

	w.Flush()
}

func recurse(level int, start, path string, parents map[string][]string) {
	// get value from map for start node
	v, ok := parents[start]
	if !ok {
		return // at a leaf node
	}

	// sort the children
	sort.Strings(v)

	level++ // increment depth
	for _, child := range v {
		looptest := *delimiter + child + *delimiter
		cycle := "No"
		if strings.Contains(path, looptest) {
			cycle = "Yes"
		}
		sLevel := fmt.Sprintf("%v", level)
		sPath := path + *delimiter + child
		writeRow(sLevel, start, child, sPath, cycle)
		if cycle == "No" {
			recurse(level, child, sPath, parents)
		}
	}

}

func writeRow(level, parent, child, path, cycle string) {
	var cells []string
	cells = append(cells, level)
	cells = append(cells, parent)
	cells = append(cells, child)
	cells = append(cells, path)
	cells = append(cells, cycle)

	err := w.Write(cells)
	if err != nil {
		log.Fatalf("csv.Write:\n%v\n", err)
	}

}

func usage(msg string) {
	fmt.Println(msg + "\n")
	flag.PrintDefaults()
	os.Exit(0)
}

var recurseHeaders []string

func init() {
	recurseHeaders = append(recurseHeaders,
		"Level", "", "", "Path", "Cycle")
}
