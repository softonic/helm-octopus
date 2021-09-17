package octopus

import (
	"os"
	"strings"
)

func GetChartPathFromArgs(args []string) (string, error) {
	path := findChartPathFromArgs(args)
	_, err := isDirectory(path)
	return path, err
}

func SwapHelmArgs(args []string, copiedFiles []CopiedFile) []string {
	helmArgs := []string{}
	for _, arg := range args {
		swapArg := ""
		for _, copiedFile := range copiedFiles {
			if arg == copiedFile.Arg {
				swapArg = copiedFile.Dst
			}
		}
		if swapArg != "" {
			helmArgs = append(helmArgs, swapArg)
		} else {
			helmArgs = append(helmArgs, arg)
		}
	}
	return helmArgs
}

func findChartPathFromArgs(args []string) string {

	i := 0
	c := 0
	skip := 0
	for {
		if i >= len(args) {
			break
		}
		arg := args[i]
		if c == 2 {
			return arg
		}
		i++
		if skip > 0 {
			skip--
			continue
		}
		// if argument is option...
		if strings.HasPrefix(arg, "-") {
			// if contains equals, then no need to skip next
			if !strings.Contains(arg, "=") {
				skip++
			}
		} else {
			c++
		}
	}
	return ""
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
