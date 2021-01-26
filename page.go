package theme

//Page the basic page fields
type Page struct {
	Title string
	Error string
}

//Context Default rendering context
var Context = struct {
	Page
	// Data map[string]interface{}
}{
	Page{Title: "My Web"},
}

//InitApp use viper to get configs
func InitApp() {

}
