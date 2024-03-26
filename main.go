package main

import (
	"os"
	"github.com/br-novelt/go-rec2csv/src/rectool"
)

func main() {
	if len(os.Args) < 2 {
		panic("No argument provided")
	}

	filename := os.Args[1]
	recfile := rectool.Load(filename)
	recfile.ToCSV()
}
