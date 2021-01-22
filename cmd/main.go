package main

import (
	"log"
	"net/http"

	"evan-soft.com/core/theme"
)

func main() {
	renderer, err := theme.Default("")
	if err != nil {
		panic(err)
	}
	//mount the theme render to root
	http.Handle("/", renderer)
	log.Println("Theme is live!")
	http.ListenAndServe(":8080", nil)
}
