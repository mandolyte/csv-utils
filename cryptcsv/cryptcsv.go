package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mandolyte/csv-utils"
)

var cs *rangespec.RangeSpec

func main() {
	e := flag.String("e", "", "Encrpytion key; required if encrypting")
	d := flag.String("d", "", "Decrpytion key; required if decrypting")
	cols := flag.String("c", "", "Range spec for columns to obfuscate")
	input := flag.String("i", "", "Input CSV filename; default STDIN")
	output := flag.String("o", "", "Output CSV filename; default STDOUT")
	headers := flag.Bool("headers", true, "CSV has headers")
	keep := flag.Bool("keep", true, "Keep CSV headers on output")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		usage("Help Message")
	}

	/* check parameters */
	if len(*e)+len(*d) == 0 {
		usage("Specify either -e or -d with key to encrypt or decrypt")
	}

	if *cols != "" {
		var cserr error
		cs, cserr = rangespec.New(*cols)
		if cserr != nil {
			log.Fatalf("Invalid column range spec:%v, Error:\n%v\n", *cols, cserr)
		}
	}

	if *keep {
		if !*headers {
			log.Fatal("Cannot keep headers you don't have!")
		}
	}

	// open output file
	var w *csv.Writer
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

	var key string
	if *e != "" {
		key = *e
	} else {
		key = *d
	}

	keydata := make([]byte, 32)
	copy(keydata, key[:])

	// read loop for CSV
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
		if (row == 0) && *headers && *keep {
			row = 1
			err := w.Write(cells)
			if err != nil {
				log.Fatalf("csv.Write:\n%v\n", err)
			}
			continue
		}
		row++
		// test columns for a match to encrypt/decrypt
		for n, v := range cells {
			if cs.InRange(uint64(n + 1)) {
				// encrpyt?
				if *d == "" {
					// decrypt key not provided so we encrypt
					cells[n] = encrypt(v, keydata)
				} else {
					// decrypt key is provided so we decrypt
					cells[n] = decrypt(v, keydata)
				}
			}
		}
		err := w.Write(cells)
		if err != nil {
			log.Fatalf("csv.Write:\n%v\n", err)
		}
	}
	w.Flush()
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: cryptcsv [options]\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func decrypt(b64 string, key []byte) string {
	// convert base64 back to byte
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Fatalf("base64 decode error:", err)
	}

	// Byte array of the string
	//ciphertext := []byte(cipherstring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("aes.NewCipher() error:", err)
	}

	// if the text is too small, then it is incorrect
	if len(data) < aes.BlockSize {
		log.Fatal("Text too short error\n")
	}

	// Get the 16 byte iv
	iv := data[:aes.BlockSize]

	// Remove it
	data = data[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(data, data)

	return string(data)
}

func encrypt(text string, key []byte) string {
	// Byte array of the string
	bytes := []byte(text)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("aes.NewCipher() error: %v\n", err)
	}

	// Create slice of (16 + bytes) length
	ciphertext := make([]byte, aes.BlockSize+len(bytes))

	// Include the IV at the beginning
	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalf("io.ReadFull() error: %v\n", err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes to ciphertext after the iv
	stream.XORKeyStream(ciphertext[aes.BlockSize:], bytes)

	// now encode to base64
	b64 := base64.StdEncoding.EncodeToString(ciphertext)
	return b64
}
