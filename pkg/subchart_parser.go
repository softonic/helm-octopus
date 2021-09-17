package octopus

import (
	"errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"path/filepath"
	"strings"
)

const subchartProtocol = "subchart://"

type SubchartParser struct {
	chart    *chart.Metadata
	basePath string
}

type Subchart struct {
	Name          string
	Alias         string
	TarFilename   string
	ArgName       string
	Id            string
	ValueFilePath string
	ChartPath     string
}

func NewSubchartParser(basePath string) (*SubchartParser, error) {
	// check if path exists
	chartYamlName := basePath + "/" + chartutil.ChartfileName
	chart, err := chartutil.LoadChartfile(chartYamlName)
	if err != nil {
		return nil, err
	}
	return &SubchartParser{
		basePath: basePath,
		chart:    chart,
	}, nil
}

func (s *SubchartParser) GetSubchartsValueFilesFromArgs(args []string) ([]Subchart, error) {
	subcharts := []Subchart{}
	for _, arg := range args {
		if strings.HasPrefix(arg, subchartProtocol) {
			s, err := s.NewSubchart(arg, s.chart.Dependencies)
			if err != nil {
				return []Subchart{}, err
			}
			subcharts = append(subcharts, s)
		}
	}
	return subcharts, nil
}

func (s *SubchartParser) NewSubchart(arg string, dependencies []*chart.Dependency) (Subchart, error) {
	subchartPath := strings.TrimPrefix(arg, subchartProtocol)
	subchartParts := strings.Split(subchartPath, "/")
	subchartId := subchartParts[0]
	subchart, err := getSubchartByID(subchartId, dependencies)
	if err != nil {
		return Subchart{}, err
	}
	subchartValueFile := filepath.Join(subchart.Name, filepath.Join(subchartParts[1:]...))
	return Subchart{
		Id:            subchartId,
		Name:          subchart.Name,
		Alias:         subchart.Alias,
		TarFilename:   getSubchartTarFilename(subchart),
		ArgName:       arg,
		ValueFilePath: subchartValueFile,
		ChartPath:     s.basePath,
	}, nil
}

func getSubchartTarFilename(subchart *chart.Dependency) string {
	return subchart.Name + "-" + subchart.Version + ".tgz"
}

func getSubchartByID(id string, dependencies []*chart.Dependency) (*chart.Dependency, error) {
	for _, dep := range dependencies {
		if id == dep.Name || id == dep.Alias {
			return dep, nil
		}
	}
	return &chart.Dependency{}, errors.New("Subchart " + id + " not found.")
}
