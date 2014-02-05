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
	db, _ := kv.CreateTemp(".", "tmp_", ".db", &kv.Options{})
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordbytes = scanner.Bytes()

		if _, db_err := db.Inc(wordbytes, 1); db_err != nil {
			log.Printf("DB ERROR: %#v", db_err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Print("scanner error: ", err)
	}

	for enum, loop_err := db.SeekFirst(); loop_err == nil; key, value, loop_err = enum.Next() {
		binary.Read(bytes.NewBuffer(value[:]), binary.BigEndian, &count)
		fmt.Printf("%v\t%d\n", string(key), count)
	}
}
