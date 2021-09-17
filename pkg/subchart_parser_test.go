package octopus

import (
	"errors"
	"helm.sh/helm/v3/pkg/chart"
	"testing"
)

func TestExtractSubchartArgs(t *testing.T) {
	args := []string{
		"upgrade",
		"my-release",
		"./my/chart/path",
		"-f",
		"subchart://foo/values.custom.yaml",
		"-f",
		"values.custom.yaml",
	}
	expected := []Subchart{
		{
			Name:        "foo",
			Alias:       "",
			TarFilename: "foo-1.2.3.tgz",
			ArgName:     "subchart://foo/values.custom.yaml",
		},
	}
	sparser := newParser()
	got, err := sparser.GetSubchartsValueFilesFromArgs(args)
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	if expected[0].Name != got[0].Name {
		t.Errorf("Expected %s, got %s", expected[0].Name, got[0].Name)
	}
	if expected[0].Alias != got[0].Alias {
		t.Errorf("Expected %s, got %s", expected[0].Alias, got[0].Alias)
	}
	if expected[0].TarFilename != got[0].TarFilename {
		t.Errorf("Expected %s, got %s", expected[0].TarFilename, got[0].TarFilename)
	}
	if expected[0].ArgName != got[0].ArgName {
		t.Errorf("Expected %s, got %s", expected[0].ArgName, got[0].ArgName)
	}
}

func TestExtractSubchartArgsByAlias(t *testing.T) {
	args := []string{
		"upgrade",
		"my-release",
		"./my/chart/path",
		"-f",
		"subchart://foobar/values.custom.yaml",
		"-f",
		"values.custom.yaml",
	}
	expected := []Subchart{
		{
			Name:        "foo",
			Alias:       "foobar",
			TarFilename: "foo-1.2.4.tgz",
			ArgName:     "subchart://foobar/values.custom.yaml",
		},
	}
	sparser := newParser()
	got, err := sparser.GetSubchartsValueFilesFromArgs(args)
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	if expected[0].Name != got[0].Name {
		t.Errorf("Expected %s, got %s", expected[0].Name, got[0].Name)
	}
	if expected[0].Alias != got[0].Alias {
		t.Errorf("Expected %s, got %s", expected[0].Alias, got[0].Alias)
	}
	if expected[0].TarFilename != got[0].TarFilename {
		t.Errorf("Expected %s, got %s", expected[0].TarFilename, got[0].TarFilename)
	}
	if expected[0].ArgName != got[0].ArgName {
		t.Errorf("Expected %s, got %s", expected[0].ArgName, got[0].ArgName)
	}
}

func TestExtractUndeclaredSubchartArgs(t *testing.T) {
	args := []string{
		"upgrade",
		"my-release",
		"./my/chart/path",
		"-f",
		"subchart://bar/values.custom.yaml",
		"-f",
		"values.custom.yaml",
	}
	expected := []Subchart{}
	expectedErr := errors.New("Subchart bar not found.")
	sparser := newParser()
	got, err := sparser.GetSubchartsValueFilesFromArgs(args)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Expected error: %v, got: %v", expectedErr, err)
	}
	if len(expected) != len(got) {
		t.Fatalf("Expected items %v, got: %v", len(expected), len(got))
	}
}

func newParser() *SubchartParser {
	deps := []*chart.Dependency{
		{
			Name:    "foo",
			Version: "1.2.3",
			Alias:   "",
		},
		{
			Name:    "foo",
			Version: "1.2.4",
			Alias:   "foobar",
		},
	}
	chart := &chart.Metadata{
		Dependencies: deps,
	}

	return &SubchartParser{chart, "mychartPath"}
}
