package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	//leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
)

var conn_leveldb *leveldb.DB

func main() {
	var err error

	ldb_dir := flag.String("l", "", "dir to save leveldb data")
	flag.Parse()

	if *ldb_dir != "" {
		conn_leveldb, err = leveldb.OpenFile(*ldb_dir, nil)
		if err != nil {
			panic(err)
		}
		defer conn_leveldb.Close()
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Specify the directory for which you desire to scan.")
		os.Exit(1)
	}

	scan_directory(args[0])
}

func scan_directory(name string) {

	files, err := ioutil.ReadDir(name)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			scan_directory(name + "/" + file.Name())
			continue
		}

		if file.Mode()&os.ModeSymlink != 0 {
			continue
		}

		md5_str, err := hash_file_md5(name + "/" + file.Name())
		if err != nil {
			panic(err)
		}

		if conn_leveldb != nil {
			if err := conn_leveldb.Put([]byte(name+"/"+file.Name()), []byte(md5_str), nil); err != nil {
				panic(err)
			}
		}

		fmt.Printf("%s  %s\n", md5_str, name+"/"+file.Name())
	}

}
