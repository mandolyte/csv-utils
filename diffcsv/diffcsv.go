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
	"time"
)

var f1name = flag.String("f1", "", "First CSV file name to compare")
var f2name = flag.String("f2", "", "Second CSV file name to compare")
var output = flag.String("o", "", "Output CSV file for differences")
var key = flag.Int("key", 0, "Key column in input CSVs (first is 1); must be unique")
var help = flag.Bool("help", false, "Show help message")
var ondupfirst = flag.Bool("ondupFirst", false, "On duplicate key, keep first one")
var onduplast = flag.Bool("ondupLast", false, "On duplicate key, keep last  one")
var noeq = flag.Bool("noeq", false, "Suppress matches, showing only differences")
var df1 = flag.String("df1", "DF1", "Alias for first input file; default DF1")
var df2 = flag.String("df2", "DF2", "Alias for second input file; default DF2")
var colnums = flag.Bool("colnums", false, "Add difference column numbers to headers")

var detailedHelp = `
	Detailed Help:
	Inputs:
		- a key column
		- two input filenames
		- an output filename
	There will be two input files to compare and there will be
	one output file created:
	a) The first file will be read and stored into a map
	b) The second file will be read and stored into a map
	c) It is an error if a file has the same key value on two rows.
	Keys must be unique within each file. 
	Note that key column number is one based, not zero based!
	NOTE! if duplicate keys exist, then there are options to keep
	the first or to keep the last one. Default is to error out.
	d) Then all keys from both inputs are combined/deduped/sorted
	e) Then we range over the combined keyset and output a new CSV
	that has a new status column as the first column and the other columns
	from the inputs as the remaining columns.
	f) the new status column has the following values:
	- EQ meaning that the values for the key are same in both input files
	- IN=1 meaning that the key and values are only in input file #1
	- IN=2 similar for input file #2
	- DFn=x,y,..,z where n is either 1 or 2; followed by a comma delimited 
	list of column numbers where the values for the key do not match.
	Note that the DF statuses always come in pairs, one for each input file.
	g) Limitations:
	- both input files must have the same number of columns
	- both must have a header row and the headers must be the same
`

func main() {
	flag.Parse()

	if *help {
		usage("")
	}

	if *key == 0 {
		usage("Key column number missing.")
	}

	if *f1name == "" {
		usage("First filename is missing.")
	}

	if *f2name == "" {
		fmt.Println()
		usage("Second filename is missing.")
	}

	if *output == "" {
		fmt.Println()
		usage("Output filename is missing.")
	}

	if *ondupfirst && *onduplast {
		fmt.Println()
		usage("Cannot use both on-dup options")
	}

	now := time.Now()
	log.Printf("Start: %v", now.Format(time.StampMilli))

	// open first input file stop.Format(Time.StampMilli)
	var r1 *csv.Reader
	f1, f1err := os.Open(*f1name)
	if f1err != nil {
		log.Fatal("os.Open() Error:" + f1err.Error())
	}
	r1 = csv.NewReader(f1)

	// open second input file
	var r2 *csv.Reader
	f2, f2err := os.Open(*f2name)
	if f2err != nil {
		log.Fatal("os.Open() Error:" + f2err.Error())
	}
	r2 = csv.NewReader(f2)

	/*********************************************************/
	// do a quick check on columns first
	// if not the same, then log error and exit

	// second file
	hdrs2, rerr := r2.Read()
	if rerr == io.EOF {
		log.Fatal("File 2 is empty", rerr)
	}
	if rerr != nil {
		log.Fatalf("csv.Read:\n%v\n", rerr)
	}
	numcols2 := len(hdrs2)

	// first file
	hdrs1, rerr := r1.Read()
	if rerr == io.EOF {
		log.Fatal("File 1 is empty", rerr)
	}
	if rerr != nil {
		log.Fatalf("csv.Read:\n%v\n", rerr)
	}
	numcols1 := len(hdrs1)

	if numcols1 != numcols2 {
		log.Fatalf("Different number of columns:%v vs. %v",
			numcols1, numcols2)
	}

	// check that headers are the same
	for i := range hdrs1 {
		if hdrs1[i] == hdrs2[i] {
			continue
		}
		log.Fatal("Headers are not the same on input files")
	}

	// check on whether to add column numbers to headers
	if *colnums {
		for i := range hdrs1 {
			hdrs1[i] = fmt.Sprintf("%v-%v", i+1, hdrs1[i])
		}
	}

	// set expectations of fields per row
	r1.FieldsPerRecord = numcols1
	r2.FieldsPerRecord = numcols1

	// open output file
	var wf1 *csv.Writer
	wf1o, wf1oerr := os.Create(*output)
	if wf1oerr != nil {
		log.Fatal("os.Create() Error:" + wf1oerr.Error())
	}
	defer wf1o.Close()
	wf1 = csv.NewWriter(wf1o)
	hdrOutput := make([]string, 0)
	hdrOutput = append(hdrOutput, "STATUS")
	hdrOutput = append(hdrOutput, hdrs1...)
	err := wf1.Write(hdrOutput)
	if err != nil {
		log.Fatalf("Output Error:\n%v\n", err)
	}

	log.Printf("Processing input #1:%v\n", *f1name)
	f1map := make(map[string][]string)
	// read first file
	rows := 0
	for {
		// read the csv file
		cells, rerr := r1.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Fatalf("csv.Read:\n%v\n", rerr)
		}
		rows++
		keyv := cells[*key-1]
		if _, ok := f1map[keyv]; ok {
			if *onduplast {
				log.Printf("Replacing non-unique key: %v on row %v\n", keyv, rows+1)
			} else if *ondupfirst {
				log.Printf("Skipping non-unique key: %v on row %v\n", keyv, rows+1)
				continue
			} else {
				log.Fatalf("Key value not unique: %v on row %v\n", keyv, rows+1)
			}
		}
		f1map[keyv] = cells
	}
	log.Printf("Number of rows in file %v:%v\n", *f1name, rows)
	f1.Close()

	log.Printf("Processing input #2:%v\n", *f2name)
	f2map := make(map[string][]string)
	// read second file
	rows = 0
	for {
		// read the csv file
		cells, rerr := r2.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Fatalf("csv.Read:\n%v\n", rerr)
		}
		rows++
		keyv := cells[*key-1]
		if _, ok := f2map[keyv]; ok {
			if *onduplast {
				log.Printf("Replacing non-unique key: %v on row %v\n", keyv, rows+1)
			} else if *ondupfirst {
				log.Printf("Skipping non-unique key: %v on row %v\n", keyv, rows+1)
				continue
			} else {
				log.Fatalf("Key value not unique: %v on row %v\n", keyv, rows+1)
			}
		}
		f2map[keyv] = cells
	}
	log.Printf("Number of rows in file %v:%v\n", *f2name, rows)
	f2.Close()

	//
	// Get a combined set of keys
	//
	uniqkeyset := make(map[string]struct{})
	for k := range f1map {
		uniqkeyset[k] = struct{}{}
	}
	for k := range f2map {
		uniqkeyset[k] = struct{}{}
	}
	keySliceSize := len(uniqkeyset)
	keys := make([]string, keySliceSize)
	slot := 0
	for k := range uniqkeyset {
		keys[slot] = k
		slot++
	}
	log.Printf("Number of combined unique keys:%v\n", keySliceSize)

	// sort them
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// counts
	eqCount := 0
	diffCount := 0
	f1UniqCount := 0
	f2UniqCount := 0

	// Now range of combined unique keys
	for n := range keys {
		val := keys[n]
		row1, ok1 := f1map[val]
		row2, ok2 := f2map[val]
		if ok1 && ok2 {
			// are all the row values the same?
			diffList := make([]int, 0)
			for i := range row1 {
				if row1[i] == row2[i] {
					continue
				}
				f := i - 1
				diffList = append(diffList, f)
			}
			if len(diffList) == 0 {
				eqCount++
				if *noeq {
					continue
				}
				outrow1 := make([]string, 0)
				outrow1 = append(outrow1, "EQ")
				outrow1 = append(outrow1, row1...)
				err := wf1.Write(outrow1)
				if err != nil {
					log.Fatalf("Output Write() Error: %v\n", err)
				}
			} else {
				diffCount++
				diffs := ""
				for i := range diffList {
					diffs += fmt.Sprintf("%v,", diffList[i]+2)
				}
				diffs = strings.TrimRight(diffs, ",")
				outrow1 := make([]string, 0)
				outrow1 = append(outrow1, fmt.Sprintf("%v=%v", *df1, diffs))
				outrow1 = append(outrow1, row1...)
				err := wf1.Write(outrow1)
				if err != nil {
					log.Fatalf("Output Write() Error: %v\n", err)
				}
				outrow2 := make([]string, 0)
				outrow2 = append(outrow2, fmt.Sprintf("%v=%v", *df2, diffs))
				outrow2 = append(outrow2, row2...)
				err = wf1.Write(outrow2)
				if err != nil {
					log.Fatalf("Output Write() Error: %v\n", err)
				}
			}
		} else {
			if !ok1 {
				f2UniqCount++
				outrow := make([]string, 0)
				outrow = append(outrow, "IN=2")
				outrow = append(outrow, row2...)
				err := wf1.Write(outrow)
				if err != nil {
					log.Fatalf("Output Write() Error: %v\n", err)
				}
			} else {
				f1UniqCount++
				outrow := make([]string, 0)
				outrow = append(outrow, "IN=1")
				outrow = append(outrow, row1...)
				err := wf1.Write(outrow)
				if err != nil {
					log.Fatalf("Output Write() Error: %v\n", err)
				}
			}
		}

	}
	wf1.Flush()

	// wrapup
	stop := time.Now()
	elapsed := time.Since(now)
	log.Printf("End: %v", stop.Format(time.StampMilli))
	log.Printf("Elapsed time %v", elapsed)

	log.Printf("------- Summary -------\n")
	log.Printf("Equal Count: %v\n", eqCount)
	log.Printf("Key Diff Count: %v\n", diffCount)
	log.Printf("Unique to input #1: %v\n", f1UniqCount)
	log.Printf("Unique to input #2: %v\n", f2UniqCount)

}

func usage(msg string) {
	fmt.Println(msg)
	fmt.Print("Usage: diffcsv [options]\n")
	flag.PrintDefaults()
	if msg == "" {
		fmt.Println(detailedHelp)
	}
	os.Exit(0)
}
