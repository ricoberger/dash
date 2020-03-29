package datasource

import (
	"errors"
	"io/ioutil"
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

type Datasource struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Auth Auth   `yaml:"auth"`
}

type Data struct {
	Labels     map[string]string
	Timestamps []int64
	Points     []float64
}

type Client interface {
	GetVariableValues(query, label string, start, end time.Time) ([]string, error)
	GetData(queries []string, start, end time.Time) ([]Data, error)
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
