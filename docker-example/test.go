package main

import "fmt"

func main() {
  a := "a"
  b := a=="a"
  fmt.Println(a)
  fmt.Println(b)
  a = "b"
  b = a=="a"
  fmt.Println(a)
  fmt.Println(b)
}
