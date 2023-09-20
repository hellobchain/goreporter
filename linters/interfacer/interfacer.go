package interfacer

import (
	"path/filepath"
	"strings"

	"github.com/hellobchain/wswlog/wlogging"
)

var logger = wlogging.MustGetLoggerWithoutName()

func Interfacer(packagesPath map[string]string) []string {
	packages := make([]string, 0)
	for _, v := range packagesPath {
		v = absPath(v)
		srcIndex := strings.Index(v, "src")
		if srcIndex >= 0 && (srcIndex+4) < len(v) {
			packages = append(packages, v[(srcIndex+4):])
		}
	}
	lines, err := CheckArgs(packages)
	if err != nil {
		logger.Error(err)
	}
	return lines
}

// absPath is a function that will get absolute path of file.
func absPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		logger.Error(err)
		return path
	}
	return absPath
}
