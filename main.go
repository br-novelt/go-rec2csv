package main

import (
    "fmt"
    "os"
)

func main() {
  fmt.Println("App starting...")
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func loadFile(filename string) {
  dat, err := os.ReadFile(filename)
  check(err)
  fmt.Print(string(dat))
}
