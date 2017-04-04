package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

var f1name *string = flag.String("f1", "", "First CSV file name to compare")
var f2name *string = flag.String("f2", "", "Second CSV file name to compare")
var help *bool = flag.Bool("help", false, "Show help message")

func main() {
	flag.Parse()

	if *help {
		usage()
	}

	if len(flag.Args()) > 0 {
		usage()
	}

	// open first input file
	var r1 *csv.Reader
	f1, f1err := os.Open(*f1name)
	if f1err != nil {
		log.Fatal("os.Open() Error:" + f1err.Error())
	}
	defer f1.Close()
	r1 = csv.NewReader(f1)

	// open second input file
	var r2 *csv.Reader
	f2, f2err := os.Open(*f2name)
	if f2err != nil {
		log.Fatal("os.Open() Error:" + f2err.Error())
	}
	defer f2.Close()
	r2 = csv.NewReader(f2)

	// ignore expectations of fields per row
	r1.FieldsPerRecord = -1
	r2.FieldsPerRecord = -1

	// open f1only file
	var wf1 *csv.Writer
	wf1o, wf1oerr := os.Create("f1only.csv")
	if wf1oerr != nil {
		log.Fatal("os.Create() Error:" + wf1oerr.Error())
	}
	defer wf1o.Close()
	wf1 = csv.NewWriter(wf1o)

	// open f2only file
	var wf2 *csv.Writer
	wf2o, wf2oerr := os.Create("f2only.csv")
	if wf2oerr != nil {
		log.Fatal("os.Create() Error:" + wf2oerr.Error())
	}
	defer wf2o.Close()
	wf2 = csv.NewWriter(wf2o)

	// open both file
	var both *csv.Writer
	botho, bothoerr := os.Create("both.csv")
	if bothoerr != nil {
		log.Fatal("os.Create() Error:" + bothoerr.Error())
	}
	defer botho.Close()
	both = csv.NewWriter(botho)

	r1csvall, r1aerr := r1.ReadAll()
	if r1aerr != nil {
		log.Fatal("r1aerr.ReadAll() Error:" + r1aerr.Error())
	}
	r2csvall, r2aerr := r2.ReadAll()
	if r2aerr != nil {
		log.Fatal("r2aerr.ReadAll() Error:" + r2aerr.Error())
	}

	// load into maps
	f1map := make(map[string]int)
	f2map := make(map[string]int)

	for n, v := range r1csvall {
		if n == 0 {
			// write header out
			err := wf1.Write(v)
			if err != nil {
				log.Fatalf("wf1.Write:\n%v\n", err)
			}
			continue
		}
	}

	for n, v := range r2csvall {
		if n == 0 {
			// write header out
			err := wf2.Write(v)
			if err != nil {
				log.Fatalf("wf2.Write:\n%v\n", err)
			}
			continue
		}
	}

}

func usage() {
	flag.PrintDefaults()
	fmt.Println("NOTE: Headers on the CSV files are expected")
	os.Exit(0)
}
