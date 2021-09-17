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

func main() {
	plugin := ""
	flag.StringVar(&plugin, "plugin", "", "Use defined plugin. Defaults to none")
	flag.Parse()
	args := flag.Args()

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
	tarHandler := octopus.NewTarHandler("/tmp/octopus/")

	copiedFiles, err := tarHandler.CreateTmpFiles(subchartValues)
	if err != nil {
		log.Fatalf("Error while creating tempfiles: %v\n", err)
	}
	helmArgs := octopus.SwapHelmArgs(args, copiedFiles)
	if plugin != "" {
		helmArgs = append([]string{plugin}, helmArgs...)
	}
	c := exec.Command("helm", helmArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	err = c.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
	}
	fmt.Fprint(os.Stdout, out.String())
	err = tarHandler.Cleanup(copiedFiles)
	if err != nil {
		log.Fatalf("Error cleaning up tempfiles: %v\n", err)
	}
}
