package main

import (
"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
