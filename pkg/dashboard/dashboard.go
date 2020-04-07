package dashboard

import (
	"errors"
	"io/ioutil"

	"github.com/ricoberger/dash/pkg/datasource"

	"gopkg.in/yaml.v2"
)

var (
	// ErrDatasourceNotFound is thrown when the datasource for the variable/graph could not be found in the provided
	// dashboard
	ErrDatasourceNotFound = errors.New("could not found datasource")
)

type Row struct {
	Height int     `yaml:"height"`
	Graphs []Graph `yaml:"graphs"`
}

type Dashboard struct {
	Name      string     `yaml:"name"`
	Variables []Variable `yaml:"variables"`
	Rows      []Row      `yaml:"rows"`
}

func New(dir string, datasources map[string]datasource.Client) ([]Dashboard, error) {
	files, err := ioutil.ReadDir(dir + "/dashboards")
	if err != nil {
		return nil, err
	}

	var dashboards []Dashboard

	for _, file := range files {
		var dashboard Dashboard

		data, err := ioutil.ReadFile(dir + "/dashboards/" + file.Name())
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &dashboard)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(dashboard.Variables); i++ {
			if client, ok := datasources[dashboard.Variables[i].Datasource]; ok {
				dashboard.Variables[i].SetClient(client)
			} else {
				return nil, ErrDatasourceNotFound
			}
		}

		for i := 0; i < len(dashboard.Rows); i++ {
			for j := 0; j < len(dashboard.Rows[i].Graphs); j++ {
				if client, ok := datasources[dashboard.Rows[i].Graphs[j].Datasource]; ok {
					dashboard.Rows[i].Graphs[j].SetClient(client)
				} else {
					return nil, ErrDatasourceNotFound
				}
			}
		}

		dashboards = append(dashboards, dashboard)
	}

	return dashboards, nil
}
