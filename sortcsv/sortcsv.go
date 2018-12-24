package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type table struct {
	records [][]string
	seq     []bool
	col     []int
}

func (t *table) Len() int {
	return len(t.records)
}

func (t *table) Swap(i, j int) {
	t.records[i], t.records[j] = t.records[j], t.records[i]
}

func (t *table) Less(i, j int) bool {
	isless := true
	for n := range t.col {
		ith := t.records[i][t.col[n]-1]
		jth := t.records[j][t.col[n]-1]
		if ith == jth {
			continue
		}
		//log.Printf("Compare %v vs %v\n", ith, jth)
		if ith < jth {
			if t.seq[n] {
				isless = true
			} else {
				isless = false
			}
		} else {
			if t.seq[n] {
				isless = false
			} else {
				isless = true
			}
		}
	}
	//log.Printf("Returning %v\n", isless)
	return isless
}

func main() {
	sortseq := flag.String("s", "", "Comma delimited list of letters 'a' or 'd', for ascending or descending (default is ascending)")
	sortcol := flag.String("c", "1", "Comma delimited list of columns to sort (default 1)")
	sortinf := flag.String("i", "", "CSV file name to sort; default STDIN")
	sortout := flag.String("o", "", "CSV output file name; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		usage()
	}

	// open output file
	var w *csv.Writer
	if *sortout == "" {
		w = csv.NewWriter(os.Stdout)
	} else {
		fo, foerr := os.Create(*sortout)
		if foerr != nil {
			log.Fatal("os.Create() Error:" + foerr.Error())
		}
		defer fo.Close()
		w = csv.NewWriter(fo)
	}

	// open input file
	var r *csv.Reader
	if *sortinf == "" {
		r = csv.NewReader(os.Stdin)
	} else {
		fi, fierr := os.Open(*sortinf)
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

	if *headers {
		werr := w.Write(csvall[0])
		if werr != nil {
			log.Fatal("w.Write() Error:" + werr.Error())
		}
		csvall = csvall[1:]
	}

	// parse columns input
	collist := strings.Split(*sortcol, ",")
	seqlist := strings.Split(*sortseq, ",")
	clist := make([]int, len(collist))
	slist := make([]bool, len(collist))
	for i := range collist {
		x, err := strconv.Atoi(collist[i])
		if err != nil {
			log.Fatalf("Element of column sort list is not an integer:%v\n", collist[i])
		}
		if x == 0 {
			log.Fatal("Column numbers begin at 1 not zero\n")
		}
		clist[i] = x
		if clist[i] > len(csvall[0]) {
			log.Fatalf("Column is larger than number of cells in row:%v\n", clist[i])
		}
		// now set the sort sequence for the column
		if i < len(seqlist) {
			if seqlist[i] == "a" || seqlist[i] == "" {
				slist[i] = true
			} else if seqlist[i] == "d" {
				slist[i] = false
			} else {
				log.Fatal("Sort sequence must 'a' for ascending or 'd' for descending\n")
			}
		} else {
			slist[i] = true
		}
	}

	/* debugging */
	/*
		log.Printf("Sort columns:%v\n", clist)
		log.Printf("Sequence columns: %v\n", slist)
	*/
	t := &table{records: csvall, seq: slist, col: clist}

	//sort.Sort(t)
	sort.Stable(t)
	werr := w.WriteAll(t.records)
	if werr != nil {
		log.Fatal("w.WriteAll() Error:" + werr.Error())
	}
	w.Flush()

}

func usage() {
	flag.PrintDefaults()
	os.Exit(0)
}
