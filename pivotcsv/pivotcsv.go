package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

var m map[string]interface{}
var hasMap bool

type table struct {
	records   [][]string
	ascending bool
	column    int
}

func (t *table) Len() int {
	return len(t.records)
}

func (t *table) Swap(i, j int) {
	t.records[i], t.records[j] = t.records[j], t.records[i]
}

func (t *table) Less(i, j int) bool {
	if hasMap {
		ival, iok := m[t.records[i][t.column]]
		if !iok {
			ival, iok = m["*"]
			if !iok {
				ival = t.records[i][t.column]
			}
		}
		jval, jok := m[t.records[j][t.column]]
		if !jok {
			jval, jok = m["*"]
			if !jok {
				jval = t.records[j][t.column]
			}
		}
		var isless bool
		switch t := ival.(type) {
		case float64:
			isless = ival.(float64) < jval.(float64)
			//log.Printf("Comparing float64 %v %v", ival, jval)
		case string:
			isless = ival.(string) < jval.(string)
			//log.Printf("Comparing STRING %v %v", ival, jval)
		default:
			log.Fatalf("Unsupported type:%T\n", t)
		}
		if t.ascending {
			return isless
		}
		return !isless
	}
	isless := t.records[i][t.column] < t.records[j][t.column]
	if t.ascending {
		return isless
	}
	return !isless
}

func main() {
	pivotcol := flag.Int("c", 0, "Column to pivot (REQUIRED)")
	pivotsum := flag.Int("s", 0, "Column to sum/concat (REQUIRED)")
	pivotinf := flag.String("i", "", "CSV file name to pivot; default STDIN")
	pivotout := flag.String("o", "", "CSV output file name; default STDOUT")
	headers := flag.Bool("headers", true, "CSV must have headers; cannot be false")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		usage()
		os.Exit(1)
	}

	if *pivotcol == 0 {
		usage()
		os.Exit(1)
	}

	if *pivotsum == 0 {
		usage()
		os.Exit(1)
	}

	if !*headers {
		usage()
		os.Exit(1)
	}

	// open output file
	var w *csv.Writer
	if *pivotout == "" {
		w = csv.NewWriter(os.Stdout)
	} else {
		fo, foerr := os.Create(*pivotout)
		if foerr != nil {
			log.Fatal("os.Create() Error:" + foerr.Error())
		}
		defer fo.Close()
		w = csv.NewWriter(fo)
		defer w.Flush()
	}

	// open input file
	var r *csv.Reader
	if *pivotinf == "" {
		r = csv.NewReader(os.Stdin)
	} else {
		fi, fierr := os.Open(*pivotinf)
		if fierr != nil {
			log.Fatal("os.Open() Error:" + fierr.Error())
		}
		defer fi.Close()
		r = csv.NewReader(fi)
	}

	// ignore expectations of fields per row
	r.FieldsPerRecord = -1

	// read into memory
	csvall, raerr := r.ReadAll()
	if raerr != nil {
		log.Fatal("r.ReadAll() Error:" + raerr.Error())
	}

	// analyze the pivot column
	// a. get list of distinct values
	// b. use this to calculate width of new CSV table
	var row int
	pivotHdrs := make(map[string]int)
	for n := range csvall {
		if row == 0 {
			row++
			continue
		}
		pivotHdrs[csvall[n][*pivotcol-1]]++
	}
	log.Printf("Number of pivot headers:%v", len(pivotHdrs))

	// sort the new pivot headers
	var phkeys []string
	for phk := range pivotHdrs {
		phkeys = append(phkeys, phk)
	}
	sort.Strings(phkeys)

	// I have enough to make the new header row now!
	// make the output slice table
	var orecs [][]string

	// let's create the header row by:
	// a. appending to a slice the non pivot and sum columns
	// b. append the phkeys from above
	var hdrrow []string
	for n, v := range csvall[0] {
		if n+1 == *pivotcol {
			// insert here the new pivot headers
			for _, w := range phkeys {
				hdrrow = append(hdrrow, w)
			}
			continue
		}
		if n+1 == *pivotsum {
			continue
		}
		hdrrow = append(hdrrow, v)
	}
	// now make the headers the first append to the table
	orecs = append(orecs, hdrrow)

	// idea:
	// create a key based on all columns EXCEPT pivot and sum
	// by letting csv package write (reduced) row to a string buffer
	// use this as a map key

	// the value for the map would be a slice of type struct:
	// type sumconcat struct {
	//   float64 -- to sum up numbers
	//   []string -- to collect non-numbers
	// }
	// the slice would be one per pivot column header value
	// or maybe a map also with header value as key and struct as value

	werr := w.WriteAll(orecs)
	if werr != nil {
		log.Fatal("w.WriteAll() Error:" + werr.Error())
	}
}

func usage() {
	flag.PrintDefaults()
	fmt.Println("tbd... other usage notes")
}
