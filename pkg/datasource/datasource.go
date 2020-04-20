package datasource

import (
	"bytes"
	"errors"
	"io/ioutil"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	// ErrInvalidType is thrown when the provided datasource in a datasource file is invalid.
	ErrInvalidType = errors.New("invalid datasource type")
)

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

type Options struct {
	MaxPoints int64 `yaml:"maxPoints"`
	Step      int64 `yaml:"step"`
}

type Datasource struct {
	Type    string  `yaml:"type"`
	Name    string  `yaml:"name"`
	URL     string  `yaml:"url"`
	Auth    Auth    `yaml:"auth"`
	Options Options `yaml:"options"`
}

type Data struct {
	Timestamps map[int]string
	Series     []Series
}

type Series struct {
	Label  string
	Points []float64
}

type TableData map[string]map[string]interface{}

type Client interface {
	GetVariableValues(query, label string, start, end time.Time) ([]string, error)
	GetData(queries, labels []string, start, end time.Time) (*Data, error)
	GetTableData(queries, labels []string) (*TableData, error)
	GetSuggestions() ([]string, error)
}

func New(dir string) (map[string]Client, error) {
	files, err := ioutil.ReadDir(dir + "/datasources")
	if err != nil {
		return nil, err
	}

	var datasources map[string]Client
	datasources = make(map[string]Client)

	for _, file := range files {
		var datasource Datasource

		data, err := ioutil.ReadFile(dir + "/datasources/" + file.Name())
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &datasource)
		if err != nil {
			return nil, err
		}

		client, err := newClient(datasource)
		if err != nil {
			return nil, err
		}

		datasources[datasource.Name] = client
	}

	return datasources, nil
}

func newClient(datasource Datasource) (Client, error) {
	switch datasource.Type {
	case "Prometheus":
		return NewPrometheusClient(datasource)
	default:
		return nil, ErrInvalidType
	}
}

func QueryInterpolation(query string, variables map[string]string) (string, error) {
	tpl, err := template.New("query").Parse(query)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, variables)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}
