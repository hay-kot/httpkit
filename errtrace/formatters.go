package errtrace

import "os"

var cwd string

// OverrideCleaner is a function that can be set to override the default path cleaner
// for the Annotate and Traceable function. The default cleaner removes the current
// working directory from the file path. Use this to override that behavior, for
// example to use the full file path.
var OverrideCleaner func(path string) string

// cleanGoPath removes the current working directory from the file path.
// This is done to make the file path relative to the current working directory.
func cleanGoPath(path string) string {
	// clean full file path to relative path
	if OverrideCleaner != nil {
		return OverrideCleaner(path)
	}

	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return path
		}
	}

	if len(path) > len(cwd) && path[:len(cwd)] == cwd {
		return path[len(cwd)+1:]
	}

	return path
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
