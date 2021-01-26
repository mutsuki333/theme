package main

import (
	"net/http"

	"github.com/mutsuki333/theme"
)

func main() {
	renderer, err := theme.Default("")
	if err != nil {
		panic(err)
	}
	//mount the theme render to root
	http.Handle("/", renderer)
	//mount the theme render to other path
	http.Handle("/page/", http.StripPrefix("/page", renderer))

	//only serve public files from a specific theme
	http.Handle("/file/", http.StripPrefix("/file", renderer.FileServer))

	//switch themes
	http.HandleFunc("/plain", func(rw http.ResponseWriter, r *http.Request) {
		renderer.Select("plain")
		rw.Write([]byte("ok"))
	})
	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		renderer.Select("test")
		rw.Write([]byte("ok"))
	})

	//another theme on different route
	renderer2 := theme.New("")
	renderer2.Select("plain")
	http.Handle("/theme2/", http.StripPrefix("/theme2", renderer))

	http.ListenAndServe(":8080", nil)
}
