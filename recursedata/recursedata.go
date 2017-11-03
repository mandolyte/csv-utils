package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var w *csv.Writer
var wpath *csv.Writer

var parent = flag.Int("parent", 1, "Parent column; default 1")
var child = flag.Int("child", 2, "Child column; default 2")
var start = flag.String("start", "", "Start value of hierarchy;\nif first letter is ampersand, use as a file with a list of values to process")
var delimiter = flag.String("delimiter", ">", "String for path delimiter")
var input = flag.String("i", "", "Input CSV filename; default STDIN")
var output = flag.String("o", "", "Output CSV filename; default STDOUT")
var headers = flag.Bool("headers", true, "Input CSV has headers")
var help = flag.Bool("help", false, "Show usage message")
var info = flag.Bool("info", true, "Show info messages during processing")
var data = flag.String("data", "", "Comma list of child data columns to include")
var pathfile = flag.String("path", "", "Output CSV file for path data")

func main() {
	flag.Parse()

	if *help {
		usage("Help Message")
	}

	if *start == "" {
		usage("Start value is missing")
	}
	var dataVals []string
	var dataVal []int
	if *data != "" {
		// split into the ints and store away for use later
		dataVals = strings.Split(*data, ",")
		dataVal = make([]int, len(dataVals))
		for i := range dataVals {
			n, err := strconv.Atoi(dataVals[i])
			if err != nil {
				log.Fatalf("strconv.Atoi() error on %v\n:%v", dataVals[i], err)
			}
			dataVal[i] = n
		}
		if *pathfile == "" {
			log.Fatal("Cannot specify data columns without a path CSV filename")
		}
	}

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

	// open output path file
	if *pathfile != "" {
		if *data == "" {
			log.Fatal("Cannot specify path CSV filename without data columns")
		}
		pfo, pfoerr := os.Create(*pathfile)
		if pfoerr != nil {
			log.Fatal("os.Create() Error:" + pfoerr.Error())
		}
		defer pfo.Close()
		wpath = csv.NewWriter(pfo)
		// write the headers
		err := wpath.Write(pathHeaders)
		if err != nil {
			log.Fatal("wpath.Write(pathHeaders) Error:" + err.Error())
		}
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

	now := time.Now().UTC()
	display(fmt.Sprintf("Start at %v", now))

	// read loop for CSV to load into memory
	var row uint64
	pcol := *parent - 1
	ccol := *child - 1
	parents := make(map[string]map[string][][]string)
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
				nil, recurseHeaders[5], recurseHeaders[6],
				true)
			row++
			continue
		}
		childmap, ok := parents[cells[pcol]]

		if ok {
			// does the child exist in the map?
			_, childOk := childmap[cells[ccol]]
			if childOk {
				// is a child table needed?
				if *data == "" {
					// no table needed
					// child is in the map already
					// nothing to do!
				} else {
					childTable := childmap[cells[ccol]]
					// child data table exists, add a new row
					newrow := make([]string, 0)
					for i := range dataVal {
						newrow = append(newrow, cells[dataVal[i]-1])
					}
					childTable = append(childTable, newrow)
					// put it back
					childmap[cells[ccol]] = childTable
				}
			} else {
				// Child is not in the map; add it
				// is a child table needed?
				if *data == "" {
					// no table needed
					childmap[cells[ccol]] = nil
				} else {
					// child data table not exists, create it first
					childTable := make([][]string, 0)
					// now make the first row for this new table
					newrow := make([]string, 0)
					for i := range dataVal {
						newrow = append(newrow, cells[dataVal[i]-1])
					}
					childTable = append(childTable, newrow)
					// put it back
					childmap[cells[ccol]] = childTable
				}
			}
			// put it back into the parent map
			parents[cells[pcol]] = childmap // do I need this??
		} else {
			// child map does not exist
			newChildMap := make(map[string][][]string)
			if *data == "" {
				// no table needed
				newChildMap[cells[ccol]] = nil
			} else {
				// child data table needed, create it first
				childTable := make([][]string, 0)
				// now make the first row for this new table
				newrow := make([]string, 0)
				for i := range dataVal {
					newrow = append(newrow, cells[dataVal[i]-1])
				}
				childTable = append(childTable, newrow)
				// put it back
				newChildMap[cells[ccol]] = childTable
			}
			// add to parent map
			parents[cells[pcol]] = newChildMap
		}
		row++
	}

	display("Data loaded and ready to start recursing")
	for _, v := range startvals {
		begin := time.Now().UTC()
		display(fmt.Sprintf("Working on %v", v))
		if *data == "" {
			recurse(0, v, v, nil, nil, parents)
		} else {
			initpath := make([]string, 0)
			initpath = append(initpath, v)
			initChildData := make([]childData, 0)
			recurse(0, v, v, initpath, initChildData, parents)
		}
		display(fmt.Sprintf(". elasped %v", time.Since(begin)))
	}
	stop := time.Now().UTC()
	elapsed := time.Since(now)
	w.Flush()
	if wpath != nil {
		wpath.Flush()
	}
	display(fmt.Sprintf("End at %v", stop))
	display(fmt.Sprintf("Elapsed time %v", elapsed))
}

type childData struct {
	child string
	data  [][]string
}

func recurse(level int, root, start string, path []string,
	pathData []childData,
	parents map[string]map[string][][]string) {

	// sort the children
	childmap := parents[start]
	var keys []string
	for k := range childmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	level++ // increment depth
	for _, child := range keys {
		cycle := contains(path, child)
		sLevel := fmt.Sprintf("%v", level)
		sPath := make([]string, len(path))
		copy(sPath, path)
		sPath = append(sPath, child)
		leaf := "Yes"
		_, ok := parents[child]
		if ok {
			leaf = "No"
		}
		writeRow(sLevel, root, start, child, sPath, leaf, cycle, false)
		var newPathData []childData
		if pathData != nil {
			newPathData = make([]childData, len(pathData))
			copy(newPathData, pathData)
			newChildData := childData{}
			newChildData.child = child
			newChildData.data = make([][]string, len(childmap[child]))
			copy(newChildData.data, childmap[child])
			newPathData = append(newPathData, newChildData)
			writePath(root, newPathData)
		}
		if cycle == "No" && ok {
			recurse(level, root, child, sPath, newPathData, parents)
		}
	}

}

func writeRow(level, root, parent, child string,
	path []string, leaf, cycle string, headerrow bool) {

	var cells []string
	cells = append(cells, level)
	cells = append(cells, root)
	cells = append(cells, parent)
	cells = append(cells, child)
	if headerrow {
		cells = append(cells, recurseHeaders[4])
	} else {
		pathString := strings.Join(path, *delimiter)
		// put a delimiter at beginning and end
		cells = append(cells, *delimiter+pathString+*delimiter)
	}
	cells = append(cells, leaf)
	cells = append(cells, cycle)

	err := w.Write(cells)
	if err != nil {
		log.Fatalf("csv.Write:\n%v\n", err)
	}
}

func writePath(root string, pathData []childData) {
	cells := make([]string, 0)
	cells = append(cells, root)
	cells = append(cells, pathData[len(pathData)-1].child)
	for _, cdata := range pathData {
		jsonVal, jsonErr := json.Marshal(cdata.data)
		if jsonErr != nil {
			log.Fatalf("json.Marshal:\n%v\n", jsonErr)
		}
		cells = append(cells, string(jsonVal))
		cells = append(cells, cdata.child)
	}
	err := wpath.Write(cells)
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
var pathHeaders []string

func init() {
	recurseHeaders = append(recurseHeaders,
		"Level", "Root", "", "", "Path", "Leaf", "Cycle")
	pathHeaders = append(pathHeaders, "root", "child",
		"data1", "child1",
		"data2", "child2",
		"data3", "child3",
		"data4", "child4",
		"data5", "child5",
		"data6", "child6",
		"data7", "child7",
		"data8", "child8",
		"data9", "child9",
		"data10", "child10",
		"data11", "child11",
		"data12", "child12",
		"data13", "child13",
		"data14", "child14",
		"data15", "child15",
	)
}

func contains(path []string, value string) string {
	for _, v := range path {
		if v == value {
			return "Yes"
		}
	}
	return "No"
}

/* Code Graveyard
func writePathRow(c []string, d []string, child string) {
	numcols := len(c) + len(d) + 1
	row := make([]string, numcols)
	i := 0
	for _, v := range c {
		row[i] = v
		i++
	}
	for _, v := range d {
		row[i] = v
		i++
	}
	row[i] = child
	err := wpath.Write(row)
	if err != nil {
		log.Fatalf("csv.Write:\n%v\n", err)
	}

}

func writePath(root string, pathData []childData) {
	cells := make([]string, 0)
	cells = append(cells, root)
	for _, cdata := range pathData {
		for _, val := range cdata.data {
			cells = append(cells, val...)
		}
		cells = append(cells, cdata.child)
	}
	err := wpath.Write(cells)
	if err != nil {
		log.Fatalf("csv.Write:\n%v\n", err)
	}
}

func writePath(root string, pathData []childData) {
	cells := make([]string, 0)
	cells = append(cells, root)
	for _, cdata := range pathData {
		jsonVal, jsonErr := json.Marshal(cdata.data)
		if jsonErr != nil {
			log.Fatalf("json.Marshal:\n%v\n", jsonErr)
		}
		cells = append(cells, string(jsonVal))
		cells = append(cells, cdata.child)
	}
	err := wpath.Write(cells)
	if err != nil {
		log.Fatalf("csv.Write:\n%v\n", err)
	}
}

*/
