package dashboard

import (
	"io/ioutil"

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

		dashboards = append(dashboards, dashboard)
	}

	return dashboards, nil
}
