package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	outputfile := flag.String("o", "", "output filename (default is stdout)")
	//dir := flag.String("d", "", "base directory to start scanning from")
	flag.Parse()

	args := flag.Args()

	fmt.Printf("%#v\n", *outputfile)
	//fmt.Printf("%#v\n", *dir)
	fmt.Printf("%#v\n", args)

	if len(args) < 1 {
		fmt.Println("Specify the directory for which you desire to scan.")
		os.Exit(1)
	}

	str, err := hash_file_md5("/etc/hosts")

	if err != nil {
		panic(err)
	}

	fmt.Println(str)
}
