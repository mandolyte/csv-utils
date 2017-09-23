package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

var help = flag.Bool("help", false, "Show help message")
var input = flag.String("input", "", "Input CSV to COPY from")
var schema = flag.String("schema", "", "Schema with table to COPY to")
var table = flag.String("table", "", "Table name to COPY to")
var urlvar = flag.String("urlvar", "",
	"Environment variable with DB URL")

func main() {
	now := time.Now()
	fmt.Printf("Start Time: %v\n", now)

	flag.Parse()
	if *help {
		usage("Help:")
	}

	if *urlvar == "" {
		usage("ERROR: Environment variable with DB URL is missing\n")
	}
	pgcred := os.Getenv(*urlvar)
	if pgcred == "" {
		usage(fmt.Sprintf("ERROR: missing URL from os.Getenv('%v')\n", pgcred))
	}

	if *schema == "" {
		usage("ERROR: Schema name with table to COPY to is missing\n")
	}

	if *table == "" {
		usage("ERROR: Table name to COPY to is missing\n")
	}

	if *input == "" {
		usage("ERROR: Input CSV to COPY from is missing\n")
	}

	db, dberr := sql.Open("postgres", pgcred)
	if dberr != nil {
		log.Fatalf("ERROR: postgres connection failed: %v", dberr)
	}

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// open input file
	var r *csv.Reader
	fi, fierr := os.Open(*input)
	if fierr != nil {
		log.Fatal("os.Open() Error:" + fierr.Error())
	}
	defer fi.Close()
	r = csv.NewReader(fi)

	// ignore expectations of fields per row
	r.FieldsPerRecord = -1

	// read the headers
	hdrs, rerr := r.Read()
	if rerr == io.EOF {
		log.Fatal("r.Read() io.EOF on first read:" + rerr.Error())
	}
	if rerr != nil {
		log.Fatal("r.Read() on headers Error:" + rerr.Error())
	}

	// ensure that all headers are lowercase
	fmt.Println("Changing all headers to lowercase!")
	for i := range hdrs {
		hdrs[i] = strings.ToLower(hdrs[i])
	}

	// prepare the statement
	stmt, err := txn.Prepare(pq.CopyInSchema(*schema, *table, hdrs...))
	if err != nil {
		log.Fatal(err)
	}

	// read loop for CSV
	var row int64
	for {
		// read the csv file
		cells, rerr := r.Read()
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Fatal("r.Read() Error:" + rerr.Error())
		}

		args := make([]interface{}, len(hdrs))
		for n := range cells {
			args[n] = &cells[n]
		}

		_, err = stmt.Exec(args...)
		if err != nil {
			log.Fatalf("Error at row %v is:\n%v\nArgs:%v", row, err, cells)
		}
		row++
	}

	// flush any remaining in buffer
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// close db
	db.Close()

	fmt.Printf("Stop Time: %v\n", time.Now())
	fmt.Printf("Total run time: %v\n", time.Since(now))
	fmt.Printf("Inserted %v rows\n", row)
}

func usage(msg string) {
	fmt.Println(msg)
	flag.PrintDefaults()
	fmt.Println("The input CSV file must have headers that match the table names")
	os.Exit(0)
}
