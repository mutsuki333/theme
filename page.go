package theme

//Page the basic page fields
type Page struct {
	Title string
	Error string
}

//Context Default rendering context
var Context = Page{
	Title: "My Theme",
}

//InitApp use viper to get configs
func InitApp() {

}
