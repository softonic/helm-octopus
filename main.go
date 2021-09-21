package main

import (
	"bytes"
	"flag" //flag "github.com/spf13/pflag"
	"fmt"
	octopus "github.com/softonic/helm-octopus/pkg"
	"log"
	"os"
	"os/exec"
)

func getEnvOrFail(env string) string {
	envVar := os.Getenv(env)
	if envVar == "" {
		log.Fatalf("Something is wrong... env var %s not defined.", env)
	}
	return envVar
}

const DEFAULT_TMP_DIR = "/tmp/octopus/"

func main() {
	helmBin := getEnvOrFail("HELM_BIN")
	pluginName := getEnvOrFail("HELM_PLUGIN_NAME")
	tmpDir := os.Getenv("HELM_OCTOPUS_TMP_DIR")

	if tmpDir == "" {
		tmpDir = DEFAULT_TMP_DIR
	}

	help := false
	plugin := ""
	flag.StringVar(&plugin, "plugin", "", "Use defined plugin. Defaults to none")
	flag.BoolVar(&help, "help", false, "Show usage")
	flag.Parse()

	if help {
		showUsage()
		os.Exit(0)
	}
	// avoid infinite recursion
	if plugin == pluginName {
		plugin = ""
	}
	args := flag.Args()

	var helmArgs []string
	if plugin != "" {
		helmArgs = []string{plugin}
	}
	if octopus.IsValidSubcommand(args) {
		tarHandler := octopus.NewTarHandler(tmpDir)

		copiedFiles := copyFilesOrFail(tarHandler, args)
		defer cleanupOrFail(tarHandler, copiedFiles)
		helmArgs = append(helmArgs, octopus.SwapHelmArgs(args, copiedFiles)...)
	} else {
		// command is invalid: pass it as is
		helmArgs = append(helmArgs, args...)
	}

	c := exec.Command(helmBin, helmArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
	}
	fmt.Fprint(os.Stdout, out.String())
}

func cleanupOrFail(tarHandler octopus.TarHandler, copiedFiles []octopus.CopiedFile) {
	err := tarHandler.Cleanup(copiedFiles)
	if err != nil {
		log.Fatalf("Error cleaning up tempfiles: %v\n", err)
	}
}

func copyFilesOrFail(tarHandler octopus.TarHandler, args []string) []octopus.CopiedFile {
	chartBasePath, err := octopus.GetChartPathFromArgs(args)
	if err != nil {
		log.Fatalf("Missing chart: %v\n", err)
	}

	subchartParser, err := octopus.NewSubchartParser(chartBasePath)
	if err != nil {
		log.Fatalf("Error while loading chart: %v\n", err)
	}

	subchartValues, err := subchartParser.GetSubchartsValueFilesFromArgs(args)
	if err != nil {
		log.Fatalf("Error while parsing args: %v\n", err)
	}

	copiedFiles, err := tarHandler.CreateTmpFiles(subchartValues)
	if err != nil {
		log.Fatalf("Error while creating tempfiles: %v\n", err)
	}

	return copiedFiles
}

func showUsage() {
	fmt.Printf(`
Usage: helm octopus <helm arguments>
  helm octopus template myrelease path/to/my/chart -f path/to/my/chart/values.yaml -f subchart://mysubchart/path/to/subchart/values.custom.yaml

  All values prefixed with subchart://, will be parsed as subchart://<subchart alias/name>/<filepath>
  
  Options:
  -help     Show this help          
  -plugin   Specify plugin to use: --plugin=<myplugin>, i.e. "helm octopus --plugin=secrets <args>" will run "helm secrets <args>"
`)
}
