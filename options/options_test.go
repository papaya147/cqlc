package options

import (
	"encoding/json"
	"testing"

	"github.com/papaya147/cqlc/util"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func createRandomOptions(t *testing.T) *Options {
	opts := Options{
		Version: int(util.RandomInt32(0, 100)),
		Cql: cqlConfig{
			Package:    util.RandomString(10),
			OutDir:     util.RandomString(10),
			CqlPackage: util.RandomString(10),
			QueriesDir: util.RandomString(10),
			SchemaDir:  util.RandomString(10),
			Overrides: []typeOverrides{
				{
					DbType: "ascii",
					GoType: util.RandomString(10),
				},
				{
					DbType: "blob",
					GoType: util.RandomString(10),
				},
			},
		},
	}

	return &opts
}

func createRandomYamlOptions(t *testing.T) *Options {
	opts := createRandomOptions(t)

	y, _ := yaml.Marshal(opts)
	err := util.WriteFile(yamlDir, y)
	require.NoError(t, err)

	return opts
}

func createRandomJsonOptions(t *testing.T) *Options {
	opts := createRandomOptions(t)

	j, _ := json.Marshal(opts)
	err := util.WriteFile(jsonDir, j)
	require.NoError(t, err)

	return opts
}

func TestLoadOptions(t *testing.T) {
	// testing yaml options
	options1 := createRandomYamlOptions(t)

	options2, err := LoadOptions()
	require.NoError(t, err)

	require.Equal(t, options1.Version, options2.Version)
	require.Equal(t, options1.Cql.Package, options2.Cql.Package)
	require.Equal(t, options1.Cql.OutDir, options2.Cql.OutDir)
	require.Equal(t, options1.Cql.CqlPackage, options2.Cql.CqlPackage)
	require.Equal(t, options1.Cql.QueriesDir, options2.Cql.QueriesDir)
	require.Equal(t, options1.Cql.SchemaDir, options2.Cql.SchemaDir)
	require.Equal(t, options1.Cql.Overrides, options2.Cql.Overrides)

	// deleting yaml file as options preferentially chooses yaml
	err = util.DeleteFile(yamlDir)
	require.NoError(t, err)

	// testing json options
	options1 = createRandomJsonOptions(t)

	options2, err = LoadOptions()
	require.NoError(t, err)

	require.Equal(t, options1.Version, options2.Version)
	require.Equal(t, options1.Cql.Package, options2.Cql.Package)
	require.Equal(t, options1.Cql.OutDir, options2.Cql.OutDir)
	require.Equal(t, options1.Cql.CqlPackage, options2.Cql.CqlPackage)
	require.Equal(t, options1.Cql.QueriesDir, options2.Cql.QueriesDir)
	require.Equal(t, options1.Cql.SchemaDir, options2.Cql.SchemaDir)
	require.Equal(t, options1.Cql.Overrides, options2.Cql.Overrides)
}
