package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	str := []byte("Some text here")
	filename := "03.txt"

	write2File(str, filename)

	readFromFile(filename)

	if err := os.Remove(filename); err != nil {
		fmt.Println("Cannot delete file", err)
	}
}

func write2File(data []byte, path2File string) {
	err := ioutil.WriteFile(path2File, data, 0600)

	if err != nil {
		fmt.Println("Cannot create file", err)

		return
	}
}

func readFromFile(path2File string) {
	data, err := ioutil.ReadFile(path2File)

	if err != nil {
		fmt.Println("Cannot open file", err)

		return
	}

	fmt.Println(data)
	fmt.Println(string(data))
}
