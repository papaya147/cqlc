package options

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type TypeOverrides struct {
	DbType string `yaml:"db_type"`
	GoType string `yaml:"go_type"`
}

type CqlConfig struct {
	Package    string          `yaml:"package"`
	OutDir     string          `yaml:"out"`
	CqlPackage string          `yaml:"cql_package"`
	Overrides  []TypeOverrides `yaml:"overrides"`
}

type Options struct {
	Version int       `yaml:"version"`
	Cql     CqlConfig `yaml:"cql"`
}

func NewOptions(path string) (*Options, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var opts Options
	if err = yaml.Unmarshal(content, &opts); err != nil {
		return nil, err
	}

	return &opts, nil
}
