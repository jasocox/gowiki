package main

import (
  "fmt"
  "net/http"
)

type Server struct {}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Hello World")
}

func main() {
  s := new(Server)
  fmt.Println("Starting Server...")

  http.ListenAndServe(":8080", s)
}
