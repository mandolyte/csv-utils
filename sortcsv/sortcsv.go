package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
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
	sortasc := flag.Bool("a", true, "Sort ascending (default true)")
	sortcol := flag.Int("c", 0, "Column to sort (default 0)")
	sortinf := flag.String("i", "", "CSV file name to sort; default STDIN")
	sortout := flag.String("o", "", "CSV output file name; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	colmap := flag.String("m", "", "Map to use instead of column values")
	flag.Parse()

	if len(flag.Args()) > 0 {
		usage()
		os.Exit(1)
	}

	//var jybte []byte
	if *colmap != "" {
		// json map available
		mi, mierr := os.Open(*colmap)
		if mierr != nil {
			log.Fatal("os.Open() Error:" + mierr.Error())
		}
		defer mi.Close()
		jbyte, jerr := ioutil.ReadAll(mi)
		if jerr != nil {
			log.Fatal("ioutil.ReadAll() Error:" + jerr.Error())
		}
		dec := json.NewDecoder(strings.NewReader(string(jbyte)))
		if err := dec.Decode(&m); err != nil {
			log.Println(err)
			return
		}
		hasMap = true
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

	t := &table{records: csvall, ascending: *sortasc, column: *sortcol}

	sort.Sort(t)
	werr := w.WriteAll(t.records)
	if werr != nil {
		log.Fatal("w.WriteAll() Error:" + werr.Error())
	}
	w.Flush()

}

func usage() {
	exampleJSON := `
{
"NP": 0,
"*": 1
}
`
	flag.PrintDefaults()
	fmt.Println("Map parameter is a filename of a JSON string")
	fmt.Println("The JSON string maps column values to values to be")
	fmt.Println("used for sorting. For example the following will")
	fmt.Println("sort all values of NP first:")
	fmt.Println(exampleJSON)
	fmt.Println("The asterisk is used to provide a default.")
	fmt.Println("A default rule is required. If the default")
	fmt.Println("mapping is also an asterisk, then the original")
	fmt.Println("value is used.")
}
