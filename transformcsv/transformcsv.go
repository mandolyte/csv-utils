package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"text/template"
	"os"
)


func main() {
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	tmplfile := flag.String("t","", "Template to use for transformation")
	output := flag.String("o", "", "Output filename; default STDOUT")
	mapname := flag.String("m", "m", "Name of map in template; default is m")
	help := flag.Bool("help", false, "Show usage message")
	flag.Parse()
 
	if *help {
		usage("Help Message")
	}

	if *tmplfile == "" {
		usage("Template file name missing")
	}
	templatebytes, terr := ioutil.ReadFile(*tmplfile)
    if terr != nil {
        log.Fatal("Template file read error:"+terr.Error())
	}
	template := string(templatebytes)


	// open output file
	var w *bufio.Writer
	if *output == "" {
		w = bufio.NewWriter(os.Stdout)
	} else {
		fo, foerr := os.Create(*output)
		if foerr != nil {
			log.Fatal("os.Create() Error:" + foerr.Error())
		}
		defer fo.Close()
		w = bufio.NewWriter(fo)
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
	r.LazyQuotes = true

	// read loop for CSV
	var hdrs []string
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
		if (row == 0) {
			row = 1
			hdrs = append(hdrs, cells...)
			continue
		}
		row++
		err := writeTemplate(w, template, *mapname, hdrs, cells)
		if err != nil {
			log.Fatal("Write error to output:"+err.Error())
		}
	}
	w.Flush()
}

func writeTemplate(w io.Writer, tmpltext,amap string, hdrs, cells []string) error {
	// logic flow 
	// 1. create a map using the hdrs as keys and cells as values
	// 2. apply the map to the template
	// 3. write it out

	m := make(map[string]string)
	for i := range hdrs {
		m[hdrs[i]] = cells [i]
	}

	t := template.Must(template.New("").Parse(tmpltext))
	err := t.Execute(w, map[string]interface{}{amap: m})
	if err != nil {
		log.Fatal("Template Execute() error:"+err.Error())
	}
	return nil
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	flag.PrintDefaults()
	os.Exit(0)
}
