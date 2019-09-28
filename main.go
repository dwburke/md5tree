package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
)

var ldb *leveldb.DB
var check_value *bool
var abs_paths *bool
var exit_val int = 0

// job workers
var WaitGroup sync.WaitGroup
var WorkerPool *workerpool.WorkerPool

func main() {
	var err error

	ldb_dir := flag.String("l", "", "dir to save leveldb data")
	check_value = flag.Bool("c", false, "check value against stored leveldb data (does not update database; '-l' is required)")
	abs_paths = flag.Bool("a", false, "resolve file paths to absolute path (i.e. ../xxx becomes /home/user/foo/xxx)")
	num_workers := flag.Int("w", 1, "number of hash workers")
	flag.Parse()

	WorkerPool = workerpool.New(*num_workers)

	ldb_dir_env := os.Getenv("MD5TREE_DATADIR")

	if *ldb_dir != "" || ldb_dir_env != "" {

		var full_path string

		if *ldb_dir != "" {
			full_path, err = filepath.Abs(*ldb_dir)
		} else {
			full_path, err = filepath.Abs(ldb_dir_env)
		}
		if err != nil {
			panic(err)
		}

		ldb, err = leveldb.OpenFile(full_path, nil)
		if err != nil {
			panic(err)
		}
		defer ldb.Close()
	}

	if ldb == nil && *check_value {
		fmt.Println("-c requires specifying -l\n")
		os.Exit(1)
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Specify the directory for which you desire to scan.")
		os.Exit(1)
	}

	path_to_scan := args[0]
	if *abs_paths {
		path_to_scan, err = filepath.Abs(path_to_scan)
		if err != nil {
			panic(err)
		}
	}
	scan_directory(path_to_scan)

	WorkerPool.Stop()
	os.Exit(exit_val)
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

		// ignore special files
		if file.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if file.Mode()&os.ModeSocket != 0 {
			continue
		}
		if file.Mode()&os.ModeDevice != 0 {
			continue
		}
		if file.Mode()&os.ModeNamedPipe != 0 {
			continue
		}
		if file.Mode()&os.ModeCharDevice != 0 {
			continue
		}

		cfile := file
		WaitGroup.Add(1)
		WorkerPool.Submit(func() {
			md5_str, err := hash_file_md5(name + "/" + cfile.Name())
			if err != nil {
				panic(err)
			}

			if ldb == nil {
				fmt.Printf("%s  %s\n", md5_str, name+"/"+cfile.Name())
			} else if ldb != nil {
				if *check_value {

					var state string = "OK"

					data, err := ldb.Get([]byte(name+"/"+cfile.Name()), nil)

					if err != nil {
						if err == leveldb_errors.ErrNotFound {
							state = "NOT FOUND"
						} else {
							panic(err)
						}
					} else if string(data) != md5_str {
						state = "FAILED"
						exit_val = 2
					}

					fmt.Printf("%s  %s: %s\n", md5_str, name+"/"+cfile.Name(), state)

				} else {
					fmt.Printf("%s  %s\n", md5_str, name+"/"+cfile.Name())

					if err := ldb.Put([]byte(name+"/"+cfile.Name()), []byte(md5_str), nil); err != nil {
						panic(err)
					}
				}

			}

			WaitGroup.Done()
		})
	}

	WaitGroup.Wait()
}
