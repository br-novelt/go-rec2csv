package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
)

type headerColumn struct {
  name string
  dataType string
  questionColumn int
  questionLine int
  questionColor int
  fieldColumn int
  fieldLine int
  fieldType string
  fieldWidth int
  entryFieldColor int
  description string
}

type header struct {
  columnsCount int
  columns []headerColumn
}

func main() {
  fmt.Println("App starting...")

  if len(os.Args) < 2 {
    panic("No argument provided")
  }

  filename := os.Args[1]
  loadFile(filename)

}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func loadFile(filename string) {
  file, err := os.Open(filename)
  check(err)

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  var text []string

  for scanner.Scan() {
      text = append(text, scanner.Text())
  }

  file.Close()

  header := newHeader(text)
  fmt.Println("**** Header loaded ****")
  fmt.Println(header)

}

/*
 * Builds a new header based on the incoming
 * string array.
 */
func newHeader(rows []string) *header {
    // Get columns count
    row := strings.Fields(rows[0])
    columnsCount, err := strconv.Atoi(row[0])
    check(err)

    h := header{columnsCount: columnsCount}

    for i := 1; i <= h.columnsCount; i++ {
        row := rows[i]
        h.columns= append(h.columns, newHeaderColumn(row))
    }

    return &h
}

/*
 * Builds a new headerColumn based on the incoming
 * string.
 */
func newHeaderColumn(row string) headerColumn {
  words := strings.Fields(row)
  hc := headerColumn{name: words[0], dataType: words[1], }
  return hc
}