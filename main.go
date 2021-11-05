package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

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
  dat, err := ioutil.ReadFile(filename)
  check(err)
  fmt.Print(string(dat))
}

/*func loadRECFile(filename string, encoding string) {
  recEof := "\r\n"
  
}*/
