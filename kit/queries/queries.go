package queries

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ReadQuery returns the query saved in a file
func ReadQuery(query string) (string, error) {
	pwd, _ := os.Getwd()
	base := strings.Split(pwd, "/test")[0]
	queryPath := filepath.Join(base, "internal/postgres/queries/", query+".sql")
	file, err := ioutil.ReadFile(queryPath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}
