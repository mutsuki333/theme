# theme

The base of core UI

## Structure

This library is an example of how to use golang's `html/template` to render sites dynamically, 
and added some personal thoughts of how to render static pages.

The file structure :

```shell
# in ./plain/
.
├── package.json        # Theme package data, the name must match the dir name
├── public              # Assets that are served by the server
│   ├── css
│   ├── img
│   └── js
└── templates           # go templates
    ├── base            # basic setup for the whole web
    │   ├── footer.tmpl
    │   ├── header.tmpl
    │   └── meta.tmpl
    ├── component       # Some other components that could be used by the engin
    │   └── xxx.tmpl
    ├── index.tmpl      # The home page when base route hit 
    ├── install.tmpl    # Every tmpl that is not in page folder are treated as standalone site,
    |                   # which do not share the property in default.tmpl,
    |                   # standalone pages are favored over pages when route hits
    ├── layout          # Other layout related tmpls
    │   ├── default.tmpl # The default layout to render
    │   └── nav
    │       └── nav.tmpl
    └── page            # folder for pages, every page is simply redefining the `content` block,
        |               # you can also override other blocks in the page tmpl,
        |               # pages will be accessible from the router. e.g. /user/profile => profile.tmpl
        ├── 404.tmpl
        ├── err.tmpl
        ├── home.tmpl
        └── user
            └── profile.tmpl
```

## Start Themeing

Install with `go get -u github.com/mutsuki333/theme/cmd/theme`.  
Copy the folder *plain* and point the path to it.

**Comands**:

```sh
theme -h
Usage: theme [options]
  -f string
        path to the theme root dir (default "./")
  -p string
        the port to listen on (default "8080")
  -s    serve index file as spa when no page found
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