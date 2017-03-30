package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var w *csv.Writer

var parent = flag.Int("parent", 1, "Parent column; default 1")
var child = flag.Int("child", 2, "Child column; default 2")
var start = flag.String("start", "", "Start value of hierarchy;\nif first letter is ampersand, use as a file with a list of values to process")
var delimiter = flag.String("delimiter", ">", "String for path delimiter")
var input = flag.String("i", "", "Input CSV filename; default STDIN")
var output = flag.String("o", "", "Output CSV filename; default STDOUT")
var headers = flag.Bool("headers", true, "Input CSV has headers")
var help = flag.Bool("help", false, "Show usage message")
var info = flag.Bool("info", true, "Show info messages during processing")

func main() {
	flag.Parse()

	if *help {
		usage("Help Message")
	}

	if *start == "" {
		usage("Start value is missing")
	}
	now := time.Now().UTC()
	display(fmt.Sprintf("Start at %v", now))

	var startvals []string
	if strings.HasPrefix(*start, "@") {
		f, ferr := os.Open((*start)[1:])
		if ferr != nil {
			log.Fatalf("os.Open() error on %v\n:%v", (*start)[1:], ferr)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			startvals = append(startvals, scanner.Text())
		}
	} else {
		startvals = append(startvals, *start)
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
			if *headers == false {
				recurseHeaders[2] = "Parent"
				recurseHeaders[3] = "Child"
			} else {
				recurseHeaders[2] = cells[pcol]
				recurseHeaders[3] = cells[ccol]
			}
			writeRow(recurseHeaders[0], recurseHeaders[1],
				recurseHeaders[2], recurseHeaders[3],
				recurseHeaders[4], recurseHeaders[5], recurseHeaders[6])
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

	display("Data loaded and ready to start recursing")
	for _, v := range startvals {
		begin := time.Now().UTC()
		display(fmt.Sprintf("Working on %v", v))
		recurse(0, v, v, *delimiter+v, parents)
		display(fmt.Sprintf(". elasped %v", time.Since(begin)))
	}
	stop := time.Now().UTC()
	elapsed := time.Since(now)
	display(fmt.Sprintf("End at %v", stop))
	display(fmt.Sprintf("Elapsed time %v", elapsed))
	w.Flush()
}

func recurse(level int, root, start, path string, parents map[string][]string) {
	// get value from map for start node
	//v, ok := parents[start]
	//if !ok {
	//	return // at a leaf node
	//}

	// sort the children
	v := parents[start]
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
		leaf := "Yes"
		_, ok := parents[child]
		if ok {
			leaf = "No"
		}
		writeRow(sLevel, root, start, child, sPath, leaf, cycle)
		if cycle == "No" && ok {
			recurse(level, root, child, sPath, parents)
		}
	}

}

func writeRow(level, root, parent, child, path, leaf, cycle string) {
	var cells []string
	cells = append(cells, level)
	cells = append(cells, root)
	cells = append(cells, parent)
	cells = append(cells, child)
	cells = append(cells, path)
	cells = append(cells, leaf)
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

func display(msg string) {
	if *info {
		log.Print(msg + "\n")
	}
}

var recurseHeaders []string

func init() {
	recurseHeaders = append(recurseHeaders,
		"Level", "Root", "", "", "Path", "Leaf", "Cycle")
}
