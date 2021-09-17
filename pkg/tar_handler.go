package octopus

import (
	"archive/tar"
	"compress/gzip"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const randomDirNameLength = 16

type Randomizer interface {
	GenerateRandomString(int) string
}

type CopiedFile struct {
	Src    string
	Dst    string
	Arg    string
	TmpDir string
}

type TarHandler struct {
	randomizer Randomizer
	basePath   string
}

type StringRandomizer struct {
}

func NewStringRandomizer(seed int64) *StringRandomizer {
	randomizer := StringRandomizer{}
	randomizer.init(seed)
	return &randomizer
}
func (sr *StringRandomizer) init(seed int64) {
	rand.Seed(seed)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (sr *StringRandomizer) GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NewTarHandler(basePath string) TarHandler {
	r := NewStringRandomizer(time.Now().UnixNano())
	return NewTarHandlerWithRandomizer(basePath, r)
}

func NewTarHandlerWithRandomizer(basePath string, randomizer Randomizer) TarHandler {
	return TarHandler{basePath: basePath, randomizer: randomizer}
}

func (t *TarHandler) CreateTmpFiles(subchartValues []Subchart) ([]CopiedFile, error) {
	copiedFiles := []CopiedFile{}
	for _, subchartValue := range subchartValues {
		c, copyErr := t.CopyTarredfile(subchartValue)
		if copyErr != nil {
			t.Cleanup(copiedFiles)
			return copiedFiles, copyErr
		}
		copiedFiles = append(copiedFiles, c)
	}
	return copiedFiles, nil
}

func (t *TarHandler) CopyTarredfile(sub Subchart) (CopiedFile, error) {
	c := CopiedFile{}
	f, err := os.Open(filepath.Join(sub.ChartPath, "/charts/", sub.TarFilename))
	if err != nil {
		return c, err
	}
	defer f.Close()
	gzf, err := gzip.NewReader(f)
	if err != nil {
		return c, err
	}
	r := tar.NewReader(gzf)
	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return c, err
		}
		if strings.HasSuffix(header.Name, sub.ValueFilePath) {
			finfo := header.FileInfo()

			randomDir := filepath.Join(t.basePath, t.randomizer.GenerateRandomString(randomDirNameLength))
			tmpFileName := filepath.Join(randomDir, sub.ValueFilePath)
			dirPath := filepath.Dir(tmpFileName)
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				return c, err
			}
			file, err := os.OpenFile(
				tmpFileName,
				os.O_RDWR|os.O_CREATE|os.O_TRUNC,
				finfo.Mode().Perm(),
			)
			if err != nil {
				return c, err
			}
			c.Src = sub.ValueFilePath
			c.Dst = tmpFileName
			c.TmpDir = randomDir
			c.Arg = sub.ArgName

			raw, err := ioutil.ReadAll(r)
			if err != nil {
				return c, nil
			}
			valueFile := map[string]interface{}{}
			err = yaml.Unmarshal(raw, valueFile)
			if err != nil {
				return c, err
			}
			valueFileWithRoot := map[string]interface{}{}
			valueFileWithRoot[sub.Id] = valueFile
			value, err := yaml.Marshal(valueFileWithRoot)
			if err != nil {
				return c, err
			}

			_, cpErr := file.Write(value)
			// _ int result is number of written bytes
			if closeErr := file.Close(); closeErr != nil {
				return c, err
			}
			if cpErr != nil {
				return c, err
			}
		}
	}
	return c, nil
}

func (t *TarHandler) Cleanup(copiedFiles []CopiedFile) error {
	for _, file := range copiedFiles {
		err := os.RemoveAll(file.TmpDir)
		if err != nil {
			return err
		}
	}
	return nil
}
