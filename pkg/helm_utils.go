package octopus

import (
	"os"
	"strings"
)

// IsValidSubcommand Checks whether subcommand is supported by octopus
func IsValidSubcommand(args []string) bool {
	var subcommands = []string{"upgrade", "install", "template", "lint"}
	firstArg := args[0]
	for _, subcommand := range subcommands {
		if firstArg == subcommand {
			return true
		}
	}
	return false
}

// GetChartPathFromArgs Finds the chart's path from the argument list and ensures it exists
func GetChartPathFromArgs(args []string) (string, error) {
	path := findChartPathFromArgs(args)
	_, err := isDirectory(path)
	return path, err
}

// SwapHelmArgs Swaps arguments matching copied file source with copied file destination
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

func isNoArgOption(opt string) bool {
	optsNoArg := []string{"--install"}
	for _, optNoArg := range optsNoArg {
		if opt == optNoArg {
			return true
		}
	}
	return false
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
			// if contains equals or is a no-argument option, then no need to skip next
			if !strings.Contains(arg, "=") && !isNoArgOption(arg) {
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
