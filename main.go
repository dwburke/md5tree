package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var output *os.File

func main() {
	outputfile := flag.String("o", "", "output filename (default is stdout)")
	//dir := flag.String("d", "", "base directory to start scanning from")
	flag.Parse()

	if *outputfile == "" {
		output = os.Stdout
	}

	fmt.Printf("%#v\n", output)
	fmt.Printf("%#v\n", os.Stdout)

	args := flag.Args()

	//fmt.Printf("%#v\n", *dir)

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

		//fmt.Printf("%#v\n", file)

		str, err := hash_file_md5(name + "/" + file.Name())

		if err != nil {
			panic(err)
		}

		fmt.Fprintf(output, "%s  %s\n", str, name+"/"+file.Name())
	}

}
