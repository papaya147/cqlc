package options

import (
	"encoding/json"
	"fmt"

	"github.com/papaya147/cqlc/util"
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

type Config struct {
	Version      int       `yaml:"version" json:"version" validate:"required"`
	Cql          cqlConfig `yaml:"cql" json:"cql" validate:"required"`
	Dependencies []string
	TypeMappings map[string]string
}

var configFileName = "cqlc"
var yamlDir = fmt.Sprintf("./%s.yaml", configFileName)
var jsonDir = fmt.Sprintf("./%s.json", configFileName)

func LoadOptions() (*Config, error) {
	configFile, fileType, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	configContent, err := util.GetFileContents(configFile)
	if err != nil {
		return nil, err
	}

	var opts Config
	switch fileType {
	case "yaml":
		if err := yaml.Unmarshal(configContent, &opts); err != nil {
			return nil, err
		}
	case "json":
		if err := json.Unmarshal(configContent, &opts); err != nil {
			return nil, err
		}
	}

	if err := util.Validate(opts); err != nil {
		return nil, err
	}

	if err := opts.replaceTypes(); err != nil {
		return nil, err
	}

	return &opts, nil
}

func getConfigFilePath() (string, string, error) {
	if err := util.CheckPathExists(yamlDir); err == nil {
		return yamlDir, "yaml", nil
	}

	if err := util.CheckPathExists(jsonDir); err == nil {
		return jsonDir, "json", nil
	}

	return "", "", fmt.Errorf("no config file found")
}
