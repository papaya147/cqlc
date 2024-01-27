package options

import (
	"fmt"
	"io"
	"os"

	"github.com/papaya147/go-cassandra-codegen/util"
	"gopkg.in/yaml.v3"
)

type typeOverrides struct {
	DbType string `yaml:"db_type" json:"db_type" validate:"required"`
	GoType string `yaml:"go_type" json:"go_type" validate:"required"`
}

type cqlConfig struct {
	Package    string          `yaml:"package" json:"package" validate:"required"`
	OutDir     string          `yaml:"out" json:"out" validate:"required"`
	CqlPackage string          `yaml:"cql_package" json:"cql_package" validate:"required"`
	QueriesDir string          `yaml:"queries" json:"queries" validate:"required"`
	SchemaDir  string          `yaml:"schema" json:"schema" validate:"required"`
	Overrides  []typeOverrides `yaml:"overrides" json:"overrides" validate:"dive"`
}

type Options struct {
	Version int       `yaml:"version" json:"version" validate:"required"`
	Cql     cqlConfig `yaml:"cql" json:"cql" validate:"required"`
}

var configFileName = "codegen"
var yamlDir = fmt.Sprintf("./%s.yaml", configFileName)
var jsonDir = fmt.Sprintf("./%s.json", configFileName)

func NewOptions() (*Options, error) {
	file, err := loadOptionsFile()
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

	if err := checkDirExists(opts.Cql.QueriesDir); err != nil {
		return nil, err
	}

	if err := checkDirExists(opts.Cql.SchemaDir); err != nil {
		return nil, err
	}

	if err := util.Validate(opts); err != nil {
		return nil, err
	}

	return &opts, nil
}

func loadOptionsFile() (*os.File, error) {
	if err := checkDirExists(yamlDir); err == nil {
		return os.Open(yamlDir)
	}

	if err := checkDirExists(jsonDir); err == nil {
		return os.Open(jsonDir)
	}

	return nil, fmt.Errorf("%s and %s config files were not found", yamlDir, jsonDir)
}

func checkDirExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("the path %s does not exist", path)
	}
	return nil
}
