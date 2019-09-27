package main

import "fmt"

func main() {
	str, err := hash_file_md5("/etc/hosts")

	if err != nil {
		panic(err)
	}

	fmt.Println(str)
}
