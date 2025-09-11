package main

import (
	"fmt"
	"os"
)

func main() {
	filename := "05.txt"

	createFile(filename)
	readFile(filename)
}

func createFile(filename string) {
	file, err := os.Create(filename)

	if err != nil {
		fmt.Println("Cannot create file.", err)

		os.Exit(1)
	}

	defer file.Close()

	data := "Some text file"

	file.WriteString(data)
}

func readFile(filename string) {
	data, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("Cannot read file", err)

		os.Exit(1)
	}

	fmt.Println(string(data))

	defer os.Remove(filename)
}
