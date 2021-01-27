package theme

//Page the basic page fields
type Page struct {
	Title string
	Error string
}

//Context Default rendering context
var Context interface{}

func init() {
	Context = struct {
		Page
		Test string
		// Data map[string]interface{}
	}{
		Page: Page{Title: "My Web"},
		Test: "cool",
	}
}

//InitApp use viper to get configs
func InitApp() {

}
