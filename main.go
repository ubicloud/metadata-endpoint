//coverage:ignore file
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	code, err := mainInternal(http.ListenAndServe)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
