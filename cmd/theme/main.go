package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mutsuki333/theme"
)

var (
	port    string
	root    string
	theming string
	SPA     bool
)

func init() {
	flag.StringVar(&port, "p", "8080", "the port to listen on")
	flag.StringVar(&root, "f", "./", "path to the theme root dir")
	flag.StringVar(&theming, "t", "plain", "theme to serve on")
	flag.BoolVar(&SPA, "s", false, "serve index file as spa when no page found")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: theme [options]\n")
	flag.PrintDefaults()
}
func main() {
	flag.Parse()
	renderer := theme.New(root)
	renderer.SPA = SPA
	err := renderer.Select(theming)
	if err != nil {
		panic(err)
	}
	//mount the theme render to root
	http.Handle("/", renderer)
	log.Println("Theme is live!")
	http.ListenAndServe(":"+port, nil)
}
