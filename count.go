package main

import (
	"bufio"
	"log"
	"github.com/cznic/kv"
	"os"
        "encoding/binary"
	"bytes"
	"fmt"
)

var (
	wordbytes []byte
	key, value []byte
	count uint64
)

func main() {
	// create database that will handle the key-value store for us
	db, _ := kv.CreateTemp(".", "tmp_", ".db", &kv.Options{})

	// close the db before the main function closes
	defer db.Close()

	// read buffered from STDIN, so that we don't have to store the whole dataset in memory
	scanner := bufio.NewScanner(os.Stdin)

	// using builtin word parser. I think it includes puncuation, so this may need to be tweaked
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		// read the parsed word from the scanner
		wordbytes = scanner.Bytes()

		// try to increment the counter in the database for the found word.
		if _, db_err := db.Inc(wordbytes, 1); db_err != nil {
			log.Printf("DB ERROR: %#v", db_err)
		}
	}

	// check if the scanner failed, and log
	if err := scanner.Err(); err != nil {
		log.Print("scanner error: ", err)
	}

	// after populating the database, let's loop through all the keys (words) and output in the required format
	for enum, loop_err := db.SeekFirst(); loop_err == nil; key, value, loop_err = enum.Next() {
		binary.Read(bytes.NewBuffer(value[:]), binary.BigEndian, &count)
		fmt.Printf("%v\t%d\n", string(key), count)
	}
}
