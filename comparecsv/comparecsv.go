package main

import (
	"crypto/sha1"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var f1name = flag.String("f1", "", "First CSV file name to compare")
var f2name = flag.String("f2", "", "Second CSV file name to compare")
var keepdups = flag.Bool("keepdups", false, "Keep duplicate rows; default is false")
var help = flag.Bool("help", false, "Show help message")

/*
	Design Overview:
	There will be two input files to compare and there will be
	three output files created:
		- f1only.csv having rows unique to f1
		- f2only.csv having rows unique to f2
		- both.csv having rows common to both f1 and f2
	a) The first file will be read and hash computed per row
	b) The hash will be a key in a map with value struct{}
	c) The second file is read and then per row:
		- the hash value is computed
		- if hash exists in first file's map, then the row is
		written to the "both" output file
		- Otherwise, it is written to the f2 only file
		- the hash is then stored similar to f1 in a map
	d) can the f1 map be reclaimed by the GC??
	e) Now f1 is read a second time and per row:
		- the hash value is computed
		- if hash exists in second file's map, then continue, since
		it is already written to the both csv file
		- otherwise, write to the f1 only file
*/

func main() {
	flag.Parse()

	if *help {
		usage()
	}

	if len(flag.Args()) > 0 {
		usage()
	}

	now := time.Now()
	log.Printf("Start at %v", now)

	// open first input file
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

	// set expectations of fields per row
	r1.FieldsPerRecord = numcols1
	r2.FieldsPerRecord = numcols1

	// open f1only file
	var wf1 *csv.Writer
	wf1o, wf1oerr := os.Create("f1only.csv")
	if wf1oerr != nil {
		log.Fatal("os.Create() Error:" + wf1oerr.Error())
	}
	defer wf1o.Close()
	wf1 = csv.NewWriter(wf1o)
	err := wf1.Write(hdrs1)
	if err != nil {
		log.Fatalf("Headers 1 Error:\n%v\n", err)
	}

	// open f2only file
	var wf2 *csv.Writer
	wf2o, wf2oerr := os.Create("f2only.csv")
	if wf2oerr != nil {
		log.Fatal("os.Create() Error:" + wf2oerr.Error())
	}
	defer wf2o.Close()
	wf2 = csv.NewWriter(wf2o)
	err = wf2.Write(hdrs2)
	if err != nil {
		log.Fatalf("Headers 2 Error:\n%v\n", err)
	}

	// open both file
	var both *csv.Writer
	botho, bothoerr := os.Create("both.csv")
	if bothoerr != nil {
		log.Fatal("os.Create() Error:" + bothoerr.Error())
	}
	defer botho.Close()
	both = csv.NewWriter(botho)
	err = both.Write(hdrs1)
	if err != nil {
		log.Fatalf("Both Headers Error:\n%v\n", err)
	}

	f1map := make(map[string]struct{})
	// read first file
	// read loop for CSV
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
		key := computeSliceSha1(cells)
		f1map[key] = struct{}{}
		rows++
	}
	log.Printf("Number of rows in file 1:%v\n", rows)
	f1.Close()

	f2map := make(map[string]struct{})
	// read second file
	// read loop for CSV
	rows = 0
	bothCount := 0
	f2Count := 0
	for {
		// read the csv file
		cells, rerr := r2.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Fatalf("csv.Read:\n%v\n", rerr)
		}
		key := computeSliceSha1(cells)
		f2map[key] = struct{}{}

		// does this row exist in file 1?
		_, f1Exists := f1map[key]
		if f1Exists {
			err := both.Write(cells)
			if err != nil {
				log.Fatalf("both Write Error:\n%v\n", err)
			}
			bothCount++
		} else {
			err := wf2.Write(cells)
			if err != nil {
				log.Fatalf("both Write Error:\n%v\n", err)
			}
			f2Count++
		}
		rows++
	}
	// flush the CSV writers
	both.Flush()
	wf2.Flush()
	f2.Close()
	botho.Close()
	log.Printf("Number of rows in file 2:%v\n", rows)
	log.Printf("Number of rows in both files:%v\n", bothCount)
	log.Printf("Number of rows ONLY in file 2:%v\n", f2Count)

	// finally re-read file 1 and match up
	// open first input file
	f1, f1err = os.Open(*f1name)
	if f1err != nil {
		log.Fatal("os.Open() Error:" + f1err.Error())
	}
	defer f1.Close()
	r1 = csv.NewReader(f1)
	f1Count := 0
	isHeader := true
	for {
		// read the csv file
		cells, rerr := r1.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Fatalf("csv.Read:\n%v\n", rerr)
		}
		if isHeader {
			isHeader = false
			continue
		}
		key := computeSliceSha1(cells)
		// does this row exist in file 2?
		_, f2Exists := f2map[key]
		if f2Exists {
			continue
		} else {
			err := wf1.Write(cells)
			if err != nil {
				log.Fatalf("both Write Error:\n%v\n", err)
			}
			f1Count++
		}
	}
	log.Printf("Number of rows ONLY in file 1:%v\n", f1Count)
	wf1.Flush()
	f1.Close()
	stop := time.Now()
	elapsed := time.Since(now)

	log.Printf("End at %v", stop)
	log.Printf("Elapsed time %v", elapsed)

}

func usage() {
	flag.PrintDefaults()
	fmt.Println("NOTE: Headers on the CSV files are expected")
	os.Exit(0)
}

func computeSliceSha1(c []string) string {
	h := sha1.New()
	for _, v := range c {
		if v == "" {
			v = "#empty"
		}
		io.WriteString(h, v)
	}
	return string(h.Sum(nil))
}

/* snippets

package main

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func main() {
	h := sha1.New()
	io.WriteString(h, "His money is twice tainted:")
	io.WriteString(h, " 'taint yours and 'taint mine.")
	fmt.Printf("% x", h.Sum(nil))
}
*/
