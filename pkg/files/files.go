package files

import (
	"path"
	"runtime"
)

func GetProjectPath() string {
	_, filename, _, _ := runtime.Caller(1)
	dir := path.Join(path.Dir(filename), "../..")

	return dir
}
