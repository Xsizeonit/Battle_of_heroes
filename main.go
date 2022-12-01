package main

import (
    "log"
    "net/http"
)

func main() {
    log.Fatal(http.ListenAndServe(":3333", http.FileServer(http.Dir("."))))
}
