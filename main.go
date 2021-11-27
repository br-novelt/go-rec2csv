package main

import (
	"fmt"
	"os"
	"rectool"
)

func main() {
	fmt.Println("App starting...")

	if len(os.Args) < 2 {
		panic("No argument provided")
	}

	filename := os.Args[1]
	recfile := rectool.Load(filename)
	recfile.ToCSV()
}
