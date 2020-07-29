package queries

import (
	"io/ioutil"
	"meli/pkg/files"
	"path/filepath"
)

// ReadQuery returns the query saved in a file
func ReadQuery(query string) (string, error) {
	projectPath := files.GetProjectPath()
	queryPath := filepath.Join(projectPath, "internal/postgres/queries/", query+".sql")
	file, err := ioutil.ReadFile(queryPath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}
