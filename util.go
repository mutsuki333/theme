package theme

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getThemeMeta(path string) (pkg map[string]interface{}, err error) {
	jsonFile, err := os.Open(filepath.Join(path, "package.json"))
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &pkg)
	return
}

func lsDir(root string) ([]string, error) {
	var files []string
	if root == "" {
		root = "."
	}
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}
