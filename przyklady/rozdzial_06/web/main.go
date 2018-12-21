package main

import (
  "flag"
  "log"
  "net/http"
)

func main() {
  var addr = flag.String("addr", ":8081", "adres witryny")
  flag.Parse()
  mux := http.NewServeMux()
  mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
  log.Println("Udostępnianie witryny na adresie:", *addr)
  http.ListenAndServe(*addr, mux)
}  