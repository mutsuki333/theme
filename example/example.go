package main

import (
	"html/template"
	"os"
)

func main() {
	var err error
	path := "./plain/templates"
	t := template.New("")
	_, err = t.ParseGlob(path + "/base/*.tmpl")
	if err != nil {
		panic(err)
	}
	_, err = t.ParseGlob(path + "/layout/*.tmpl")
	if err != nil {
		panic(err)
	}
	_, err = t.ParseGlob(path + "/component/*.tmpl")
	if err != nil {
		panic(err)
	}
	_, err = t.ParseGlob(path + "/index.tmpl")
	if err != nil {
		panic(err)
	}

	n, err := t.Clone()
	if err != nil {
		panic(err)
	}
	_, err = n.ParseFiles(path + "/page/home.tmpl")
	if err != nil {
		panic(err)
	}
	err = n.ExecuteTemplate(os.Stdout, "index", map[string]interface{}{"Title": "Hello"})
	if err != nil {
		panic(err)
	}

	n, err = t.Clone()
	if err != nil {
		panic(err)
	}
	_, err = n.ParseFiles(path + "/page/install.tmpl")
	if err != nil {
		panic(err)
	}
	err = n.ExecuteTemplate(os.Stdout, "install", map[string]interface{}{"Title": "Hello"})
	if err != nil {
		panic(err)
	}

	n, err = t.Clone()
	if err != nil {
		panic(err)
	}
	_, err = n.ParseFiles(path + "/page/404.tmpl")
	if err != nil {
		panic(err)
	}
	err = n.ExecuteTemplate(os.Stdout, "index", map[string]interface{}{"Title": "Hello"})
	if err != nil {
		panic(err)
	}
}
