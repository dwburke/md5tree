package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
)

var ldb *leveldb.DB
var compare_before_save *bool

func main() {
	var err error

	ldb_dir := flag.String("l", "", "dir to save leveldb data")
	compare_before_save = flag.Bool("c", false, "compare before saving to leveldb")
	flag.Parse()

	if *ldb_dir != "" {
		ldb, err = leveldb.OpenFile(*ldb_dir, nil)
		if err != nil {
			panic(err)
		}
		defer ldb.Close()
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

		if ldb == nil {
			fmt.Printf("%s  %s\n", md5_str, name+"/"+file.Name())
		} else if ldb != nil {
			if *compare_before_save {

				var state string = "OK"

				data, err := ldb.Get([]byte(name+"/"+file.Name()), nil)

				if err != nil {
					if err == leveldb_errors.ErrNotFound {
						state = "NOT FOUND"
					} else {
						panic(err)
					}

				} else if string(data) != md5_str {
					state = "FAILED"
				}

				fmt.Printf("%s  %s: %s\n", md5_str, name+"/"+file.Name(), state)
			} else {
				fmt.Printf("%s  %s\n", md5_str, name+"/"+file.Name())
			}

			if err := ldb.Put([]byte(name+"/"+file.Name()), []byte(md5_str), nil); err != nil {
				panic(err)
			}
		}

	}

}
