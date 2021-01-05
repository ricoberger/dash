package dashboard

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Row struct {
	Height int     `yaml:"height"`
	Graphs []Graph `yaml:"graphs"`
}

type Dashboard struct {
	Name              string     `yaml:"name"`
	DefaultDatasource string     `yaml:"defaultDatasource"`
	Variables         []Variable `yaml:"variables"`
	Rows              []Row      `yaml:"rows"`
}

func New(dir string) ([]Dashboard, error) {
	dashboardDir := filepath.Join(dir, "dashboards")
	
	files, err := ioutil.ReadDir(dashboardDir)
	if err != nil {
		return nil, err
	}

	var dashboards []Dashboard

	for _, file := range files {
		var dashboard Dashboard

		dashboardFile := filepath.Join(dashboardDir, file.Name())

		data, err := ioutil.ReadFile(dashboardFile)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &dashboard)
		if err != nil {
			return nil, err
		}

		dashboards = append(dashboards, dashboard)
	}

	return dashboards, nil
}

func Explore(query string) ([]Dashboard, error) {
	dashboard := Dashboard{
		Name: "Explore",
		Rows: []Row{
			{
				Height: 99,
				Graphs: []Graph{
					{
						Width: 99,
						Type:  "linechart",
						Title: "Explore",
						Queries: []Query{
							{
								Query: query,
							},
						},
						Options: Options{
							Legend: "bottom",
						},
					},
				},
			},
		},
	}

	return []Dashboard{dashboard}, nil
}
