package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"evan-soft.com/core/theme"
)

var (
	port    string
	root    string
	theming string
)

func init() {
	flag.StringVar(&port, "p", "8080", "the port to listen on")
	flag.StringVar(&root, "f", "./", "path to the theme root dir")
	flag.StringVar(&theming, "t", "plain", "theme to serve on")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: theme [options]\n")
	flag.PrintDefaults()
}
func main() {
	flag.Parse()
	renderer := theme.New(root)
	err := renderer.Select(theming)
	if err != nil {
		panic(err)
	}
	//mount the theme render to root
	http.Handle("/", renderer)
	log.Println("Theme is live!")
	http.ListenAndServe(":"+port, nil)
}
