package theme

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	FileType   string
	Layout     string
	Home       string
	SPA        bool
	Debug      bool
}

//Default plain theme
func Default(root string) (*Renderer, error) {
	r := New(root)
	err := r.Select("plain")
	return r, err
}

//New Renderer
func New(root string) *Renderer {
	r := &Renderer{
		Root:     root,
		FileType: ".tmpl",
		Layout:   "default",
		Home:     "index",
		SPA:      false,
		Debug:    true,
	}
	return r
}

//Themes list available themes
func (r *Renderer) Themes() []string {
	var themes []string
	dirs, err := lsDir(r.Root)
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
			if !info.IsDir() && strings.HasSuffix(path, r.FileType) {
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
	r.FileServer = http.FileServer(http.Dir(filepath.Join(r.Root, r.Theme, "public")))
	return err
}

//RenderPage render
func (r *Renderer) RenderPage(wr io.Writer, path string, data interface{}, layout ...string) error {
	var err error
	page, err := r.tpl.Clone()
	if err != nil {
		return err
	}
	_, err = page.ParseFiles(filepath.Join(r.Root, r.Theme, "templates", "page", path+r.FileType))
	if err != nil {
		return err
	}
	lay := r.Layout
	if len(layout) > 0 {
		lay = layout[0]
	}
	buf := &bytes.Buffer{}
	err = page.ExecuteTemplate(buf, lay, data)
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
	_, err = page.ParseFiles(filepath.Join(r.Root, r.Theme, "templates", path+r.FileType))
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
	path := strings.TrimPrefix(req.URL.Path, "/")
	var err error
	if path == "" {
		err = r.RenderStandalonePage(w, r.Home, Context)
	} else if r.HasSPage(path) {
		err = r.RenderStandalonePage(w, path, Context)
	} else if r.HasPath(path) {
		err = r.RenderPage(w, path, Context)
	} else if r.HasAsset(path) {
		r.FileServer.ServeHTTP(w, req)
	} else if r.SPA {
		err = r.RenderStandalonePage(w, r.Home, Context)
	} else {
		ctx := Context
		w.WriteHeader(http.StatusNotFound)
		r.RenderPage(w, "404", ctx)
	}
	if err != nil {
		ctx := Page{}
		ctx.Error = err.Error()
		log.Println(err)
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

//HasSPage file or path in theme folder
func (r *Renderer) HasSPage(path string) bool {
	stats, err := os.Stat(filepath.Join(r.tplDir, path+r.FileType))
	if err != nil || stats.IsDir() {
		return false
	}
	return true
}

//HasPath file or path in theme folder
func (r *Renderer) HasPath(path string) bool {
	stats, err := os.Stat(filepath.Join(r.tplDir, "page", path+r.FileType))
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
