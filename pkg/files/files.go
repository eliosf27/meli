package files

import (
	"path"
	"runtime"
)

// GetProjectPath returns the root dir of the project
func GetProjectPath() string {
	_, filename, _, _ := runtime.Caller(1)
	dir := path.Join(path.Dir(filename), "../..")

	return dir
}
