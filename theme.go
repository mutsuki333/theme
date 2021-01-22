package theme

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	FileType = ".tmpl"
	Entry    = "entry"
	Index    = "home"
)

//ModuleDirs path to look for template module files
var ModuleDirs = []string{"base", "layout", "component"}

//Renderer that parse
type Renderer struct {
	Root       string
	Theme      string
	Pkg        map[string]interface{}
	tpl        *template.Template
	tplDir     string
	FileServer http.Handler
}

//Default plain theme
func Default(root string) (*Renderer, error) {
	r := &Renderer{
		Root: root,
	}
	err := r.Select("plain")
	return r, err
}

//New Renderer
func New(root string) *Renderer {
	r := &Renderer{
		Root: root,
	}
	return r
}

//Themes list available themes
func (r *Renderer) Themes() []string {
	var themes []string
	dirs, err := lsDir(r.Root)
	fmt.Println(dirs)
	if err != nil {
		return themes
	}
	for _, dir := range dirs {
		theme, err := getThemeMeta(filepath.Join(r.Root, dir))
		nameI, ok := theme["name"]
		if err != nil || !ok {
			continue
		}
		name, ok := nameI.(string)
		if ok {
			themes = append(themes, name)
		}
	}
	return themes
}

//Template get the clone of the template
func (r *Renderer) Template() *template.Template {
	n, err := r.tpl.Clone()
	if err != nil {
		panic(err)
	}
	return n
}

//Add a template
func (r *Renderer) Add(text string) error {
	_, err := r.tpl.Parse(text)
	return err
}

//Select a theme
func (r *Renderer) Select(theme string) error {
	var err error
	r.Pkg, err = getThemeMeta(filepath.Join(r.Root, theme))
	_, ok := r.Pkg["name"]
	if err != nil || !ok {
		return err
	}
	r.Theme = theme
	r.tplDir = filepath.Join(r.Root, theme, "templates")
	r.tpl = template.New("")
	for _, dir := range ModuleDirs {
		err = filepath.Walk(filepath.Join(r.tplDir, dir), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, FileType) {
				_, err = r.tpl.ParseFiles(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	_, err = r.tpl.ParseFiles(filepath.Join(r.tplDir, Entry+FileType))
	r.FileServer = http.FileServer(http.Dir(filepath.Join(r.Root, r.Theme, "public")))
	return err
}

//RenderPage render
func (r *Renderer) RenderPage(wr io.Writer, path string, data interface{}) error {
	var err error
	page, err := r.tpl.Clone()
	if err != nil {
		return err
	}
	_, err = page.ParseFiles(filepath.Join(r.Root, r.Theme, "templates", "page", path+FileType))
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = page.ExecuteTemplate(buf, Entry, data)
	if err != nil {
		return err
	}
	buf.WriteTo(wr)
	return nil
}

//RenderStandalonePage render
func (r *Renderer) RenderStandalonePage(wr io.Writer, path string, data interface{}) error {
	var err error
	if strings.Contains(path, "/") {
		return errors.New("standalone pages should not be in sub folder")
	}
	page, err := r.tpl.Clone()
	if err != nil {
		return err
	}
	_, err = page.ParseFiles(filepath.Join(r.Root, r.Theme, "templates", path+FileType))
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = page.ExecuteTemplate(buf, path, data)
	if err != nil {
		return err
	}
	buf.WriteTo(wr)
	return nil
}

//ServeHTTP serve all file and templates, 404 if notfound
func (r *Renderer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	var err error
	if path == "/" {
		err = r.RenderPage(w, Index, Context)
	} else if r.HasPath(path) {
		err = r.RenderPage(w, path, Context)
	} else if r.HasAsset(path) {
		r.FileServer.ServeHTTP(w, req)
	} else {
		ctx := Context
		w.WriteHeader((http.StatusNotFound))
		r.RenderPage(w, "404", ctx)
	}
	if err != nil {
		ctx := Context
		ctx.Error = err.Error()
		r.RenderPage(w, "err", ctx)
	}

}

//ServeFS file and templates
func (r *Renderer) ServeFS(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if r.HasAsset(path) {
		r.FileServer.ServeHTTP(w, req)
		return
	}
	w.WriteHeader((http.StatusNotFound))

}

//HasPath file or path in theme folder
func (r *Renderer) HasPath(path string) bool {
	stats, err := os.Stat(filepath.Join(r.tplDir, "page", path+FileType))
	if err != nil || stats.IsDir() {
		return false
	}
	return true
}

//HasAsset has asset in theme public folder
func (r *Renderer) HasAsset(file string) bool {
	stats, err := os.Stat(filepath.Join(r.Root, r.Theme, "public", file))
	if err != nil || stats.IsDir() {
		return false
	}
	return true
}
