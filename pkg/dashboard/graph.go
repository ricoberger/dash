package dashboard

import (
	"time"

	"github.com/ricoberger/dash/pkg/datasource"
)

type Graph struct {
	client datasource.Client `yaml:"-"`

	Width      float64  `yaml:"width"`
	Datasource string   `yaml:"datasource"`
	Type       string   `yaml:"type"`
	Title      string   `yaml:"title"`
	Queries    []string `yaml:"queries"`
	Options    Options  `yaml:"options"`
}

type Options struct {
	Unit       string    `yaml:"unit"`
	Stats      []string  `yaml:"stats"`
	Prefix     string    `yaml:"prefix"`
	Postfix    string    `yaml:"postfix"`
	Decimals   int       `yaml:"decimals"`
	Thresholds []float64 `yaml:"thresholds"`
	Colors     []string  `yaml:"colors"`
	Label      string    `yaml:"label"`
}

func (g *Graph) SetClient(client datasource.Client) {
	g.client = client
}

func (g *Graph) GetData(variables map[string]string, start, end time.Time) ([]datasource.Data, error) {
	var queries []string

	for _, query := range g.Queries {
		q, err := queryInterpolation(query, variables)
		if err != nil {
			return nil, err
		}

		queries = append(queries, q)
	}

	return g.client.GetData(queries, start, end)
}
