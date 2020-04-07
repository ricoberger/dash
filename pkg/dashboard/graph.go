package dashboard

import (
	"time"

	"github.com/ricoberger/dash/pkg/datasource"
)

type Graph struct {
	client datasource.Client `yaml:"-"`

	Width      int     `yaml:"width"`
	Datasource string  `yaml:"datasource"`
	Type       string  `yaml:"type"`
	Title      string  `yaml:"title"`
	Queries    []Query `yaml:"queries"`
	Options    Options `yaml:"options"`
}

type Query struct {
	Query string `yaml:"query"`
	Label string `yaml:"label"`
}

type Options struct {
	Unit       string    `yaml:"unit"`
	Stats      []string  `yaml:"stats"`
	Decimals   int       `yaml:"decimals"`
	Thresholds []float64 `yaml:"thresholds"`
	Colors     []string  `yaml:"colors"`
	Legend     string    `yaml:"legend"`
}

func (g *Graph) SetClient(client datasource.Client) {
	g.client = client
}

func (g *Graph) GetData(variables map[string]string, start, end time.Time) (*datasource.Data, error) {
	var queries []string
	var labels []string

	for _, query := range g.Queries {
		q, err := datasource.QueryInterpolation(query.Query, variables)
		if err != nil {
			return nil, err
		}

		queries = append(queries, q)
		labels = append(labels, query.Label)
	}

	return g.client.GetData(queries, labels, start, end)
}
