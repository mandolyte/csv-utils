package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type table struct {
	records                 [][]string
	pivotcol, pivotcolcount int
}

func (t *table) Len() int {
	return len(t.records)
}

func (t *table) Swap(i, j int) {
	t.records[i], t.records[j] = t.records[j], t.records[i]
}

func (t *table) Less(i, j int) bool {
	for n := range t.records[i] {
		if n >= t.pivotcol && n < (t.pivotcol+t.pivotcolcount) {
			continue
		}
		if t.records[i][n] < t.records[j][n] {
			return true
		}
	}
	return false
}

func main() {
	pivotcol := flag.Int("c", 0, "Column to pivot (REQUIRED)")
	pivotsum := flag.Int("s", 0, "Column to sum/concat (REQUIRED)")
	pivotinf := flag.String("i", "", "CSV file name to pivot; default STDIN")
	pivotout := flag.String("o", "", "CSV output file name; default STDOUT")
	headers := flag.Bool("headers", true, "CSV must have headers; cannot be false")
	help := flag.Bool("help", false, "Show help message")
	novalue := flag.String("nv", "", "String to signal novalue; default is empty string")
	numformat := flag.String("nf", "%v", "Format to use for numbers")
	onlynum := flag.Bool("on", true, "Only consider numeric data and sum them")
	onlystr := flag.Bool("os", false, "Consider data as strings and concatenate")
	strdlm := flag.String("sd", ",", "Concatenation delimiter; default is comma")
	flag.Parse()

	if *help {
		usage("")
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		usage("Arguments provided when none expected")
		os.Exit(1)
	}

	if *pivotcol == 0 {
		usage("Pivot column number must greater than zero")
		os.Exit(0)
	}

	if *pivotsum == 0 {
		usage("Pivot sum column number must greater than zero")
		os.Exit(0)
	}

	if !*headers {
		usage("Headers are required; add them before using")
		os.Exit(0)
	}

	if *onlystr {
		*onlynum = false
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
	//log.Printf("Number of pivot headers:%v", len(pivotHdrs))

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
	type sumconcat struct {
		sumnum float64
		sumstr []string
		ncount uint64
	}
	pivot := make(map[string](map[string]*sumconcat))
	for _, v := range csvall[1:] {
		//
		// step 1. create the []string for the key
		//
		var tmp []string
		for x, y := range v {
			// skip the pivot and sum columns
			if x+1 == *pivotcol {
				continue
			}
			if x+1 == *pivotsum {
				continue
			}
			tmp = append(tmp, y)
		}
		//
		// step 2. let CSV package create the key
		//
		var b bytes.Buffer
		w2 := csv.NewWriter(&b)
		err := w2.Write(tmp)
		w2.Flush()
		if err != nil {
			log.Fatal("w2.Write() Error:" + err.Error())
		}
		skey := b.String()

		// step 3. update key value
		mapsc, ok := pivot[skey]
		//fmt.Printf("Summing:%v\n", v[*pivotsum-1])
		//fmt.Printf("Pivot col value is:%v\n", v[*pivotcol-1])
		if ok {
			// if key exists already in the pivot map
			// update the values and continue
			// try to convert pivotsum column value to a float64
			/*
				fmt.Print("Working on map:\n")

				for debugk, debugv := range mapsc {
					fmt.Printf("Key: %v -- Val: %v\n", debugk, debugv)
				}
			*/
			if *onlynum {
				if f, err := strconv.ParseFloat(v[*pivotsum-1], 64); err == nil {
					tmpsc, ok := mapsc[v[*pivotcol-1]]
					if ok {
						mapsc[v[*pivotcol-1]].sumnum += f
						mapsc[v[*pivotcol-1]].ncount++
					} else {
						tmpsc = new(sumconcat)
						tmpsc.sumnum = f
						tmpsc.ncount++
						mapsc[v[*pivotcol-1]] = tmpsc
					}
				}
			} else {
				tmpsc, ok := mapsc[v[*pivotcol-1]]
				if ok {
					mapsc[v[*pivotcol-1]].sumstr =
						append(mapsc[v[*pivotcol-1]].sumstr, v[*pivotsum-1])
				} else {
					tmpsc = new(sumconcat)
					tmpsc.sumstr = append(tmpsc.sumstr, v[*pivotsum-1])
					mapsc[v[*pivotcol-1]] = tmpsc
				}
			}
		} else {
			//
			// step 3b. fill out the struct val for map
			//
			sc := new(sumconcat)
			if *onlynum {
				if f, err := strconv.ParseFloat(v[*pivotsum-1], 64); err == nil {
					sc.sumnum = f
				}
			} else {
				sc.sumstr = append(sc.sumstr, v[*pivotsum-1])
			}
			tmpmap := make(map[string]*sumconcat)
			tmpmap[v[*pivotcol-1]] = sc
			pivot[skey] = tmpmap
		}

	}
	csvall = nil
	// now create the output table
	for k, v := range pivot {
		// untangle the CSV formatted key using CSV package
		//fmt.Printf("Pivot Key is /%v/\n", k)
		b := bytes.NewBufferString(k)
		r := csv.NewReader(b)
		row, rerr := r.Read()
		if rerr != nil {
			if rerr != io.EOF {
				log.Fatal("r.Read Error:" + rerr.Error())
			}
		}
		// append to a new row, inserting pivot columns in correct spot
		var newrow []string
		for i := 0; i < *pivotcol-1; i++ {
			if i == (*pivotsum - 1) {
				continue
			}
			newrow = append(newrow, row[i])
		}
		// now for the pivot columns
		// use the sorted slice to pick them in sorted order
		for _, vsc := range phkeys {
			//fmt.Printf("phkey is /%v/\n", vsc)
			sc, ok := v[vsc]
			if ok {
				//fmt.Printf("Found v[vsc] /%v/\n", sc)
				if *onlynum {
					if sc.ncount == 0 {
						newrow = append(newrow, *novalue)
					} else {
						newrow = append(newrow, fmt.Sprintf(*numformat, sc.sumnum))
					}
				} else {
					newrow = append(newrow, strings.Join(sc.sumstr, *strdlm))
				}
			} else {
				//fmt.Printf("NOT Found v[vsc] /%v/\n", sc)
				// nothing for this header key - put out empty strings
				newrow = append(newrow, *novalue)
			}
		}
		// now append the rest of the columns after *pivotcol
		for i := *pivotcol; i < len(row); i++ {
			if i == (*pivotsum - 1) {
				continue
			}
			newrow = append(newrow, row[i])
		}

		// append row to orecs table
		orecs = append(orecs, newrow)
	}

	// write out the header row
	werr := w.Write(orecs[0])
	if werr != nil {
		log.Fatal("w.Write() Error:" + werr.Error())
	}

	// now, let's sort the table
	tbl := &table{records: orecs[1:], pivotcol: *pivotcol - 1, pivotcolcount: len(phkeys)}

	sort.Sort(tbl)

	werr = w.WriteAll(tbl.records)
	if werr != nil {
		log.Fatal("w.WriteAll() Error:" + werr.Error())
	}
}

func usage(msg string) {
	if msg != "" {
		fmt.Println(msg)
	}
	flag.PrintDefaults()
}
