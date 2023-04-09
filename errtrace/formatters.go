package errtrace

import (
	"os"
	"strings"
)

// OverrideCleaner is a function that can be set to override the default path cleaner
// for the Annotate and Traceable function. The default cleaner removes the current
// working directory from the file path. Use this to override that behavior, for
// example to use the full file path.
var OverrideCleaner func(path string) string

func cleanGoPath(path string) string {
	// clean full file path to relative path
	if OverrideCleaner != nil {
		return OverrideCleaner(path)
	}

	return path
}

// RelativeCleaner returns a function that can be used to clean the file path
// to a relative path. The returned function will remove the current working
// directory from the file path.
//
// If the current working directory cannot be determined, the returned function
// will fall back to the default path cleaner.
//
// Example:
//
//	errtrace.OverrideCleaner = errtrace.RelativeCleaner()
func RelativeCleaner() func(path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return cleanGoPath
	}

	return func(path string) string {
		if strings.HasPrefix(path, cwd) {
			return path[len(cwd)+1:]
		}
		return cleanGoPath(path)
	}
}

func trimFuncName(name string) string {
	lastPeriod := 0

	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			lastPeriod = i
		}
		if name[i] == '/' {
			return name[lastPeriod+1:]
		}
	}
	return name
}
