# theme

The base of core UI

## Structure

This library is a example of how to use golang's `html/template` to render sites dynamically, the file struture

```
│  package.json
│
├─public
│  ├─css
│  ├─img
│  └─js
└─templates
	│ ├─ entry.tmpl
	│ └─ install.tmpl
	├─base
	│ ├─ footer.tmpl
	│ ├─ header.tmpl
	│ └─ meta.tmpl
	├─component
	│ └─ hello_world.tmpl
	├─layout
	│  ├─ layout.tmpl
	│  └─nav
	│     └─ nav.tmpl
	└─page
		├─ 404.tmpl
		├─ err.tmpl
		├─ home.tmpl
		└─user
			└─ profile.tmpl
```


## Start Themeing

fork this repo and run

```shell
make run -p PORT
```

**Comands**:

```sh
theme -h
Usage: theme [options]
  -f string
        path to the theme root dir (default "./")
  -p string
        the port to listen on (default "8080")
  -t string
        theme to serve on (default "plain")
```


## To use and extend the library

```go

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

```