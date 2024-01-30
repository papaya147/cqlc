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
	Version      int       `yaml:"version" json:"version" validate:"required"`
	Cql          cqlConfig `yaml:"cql" json:"cql" validate:"required"`
	Dependencies []string
	TypeMappings map[string]string
}

var configFileName = "codegen"
var yamlDir = fmt.Sprintf("./%s.yaml", configFileName)
var jsonDir = fmt.Sprintf("./%s.json", configFileName)

func NewOptions() (*Options, error) {
	var opts Options

	file, err := opts.loadConfigFile()
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(content, &opts); err != nil {
		return nil, err
	}

	if err := util.CheckPathExists(opts.Cql.QueriesDir); err != nil {
		return nil, err
	}

	if err := util.CheckPathExists(opts.Cql.SchemaDir); err != nil {
		return nil, err
	}

	if err := util.Validate(opts); err != nil {
		return nil, err
	}

	if err := opts.replaceTypes(); err != nil {
		return nil, err
	}

	return &opts, nil
}

func (opts *Options) loadConfigFile() (*os.File, error) {
	file, err := util.GetFile(yamlDir)
	if err == nil {
		return file, nil
	}

	file, err = util.GetFile(jsonDir)
	if err == nil {
		return file, nil
	}

	return nil, fmt.Errorf("%s and %s config files were not found", yamlDir, jsonDir)
}
