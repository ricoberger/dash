package datasource

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Prometheus struct {
	v1api v1.API
}

type basicAuthTransport struct {
	Transport http.RoundTripper

	username string
	password string
}

type tokenAuthTransporter struct {
	Transport http.RoundTripper

	token string
}

func (bat basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(bat.username, bat.password)
	return bat.Transport.RoundTrip(req)
}

func (tat tokenAuthTransporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+tat.token)
	return tat.Transport.RoundTrip(req)
}

func NewPrometheusClient(datasource Datasource) (*Prometheus, error) {
	roundTripper := api.DefaultRoundTripper

	if datasource.Auth.Username != "" && datasource.Auth.Password != "" {
		roundTripper = basicAuthTransport{
			Transport: roundTripper,
			username:  datasource.Auth.Username,
			password:  datasource.Auth.Password,
		}
	}

	if datasource.Auth.Token != "" {
		roundTripper = tokenAuthTransporter{
			Transport: roundTripper,
			token:     datasource.Auth.Token,
		}
	}

	client, err := api.NewClient(api.Config{
		Address:      datasource.URL,
		RoundTripper: roundTripper,
	})

	if err != nil {
		return nil, err
	}

	return &Prometheus{
		v1api: v1.NewAPI(client),
	}, nil
}

func (p *Prometheus) GetVariableValues(query, label string, start, end time.Time) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	labelSets, _, err := p.v1api.Series(ctx, []string{query}, start, end)
	if err != nil {
		return nil, err
	}

	var values []string

	for _, labelSet := range labelSets {
		for key, value := range labelSet {
			if string(key) == label {
				values = append(values, string(value))
			}
		}
	}

	return values, nil
}

func (p *Prometheus) GetData(queries []string, start, end time.Time) ([]Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var queriesData []Data

	for _, query := range queries {
		timeRange := v1.Range{
			Start: start,
			End:   end,
			Step:  10 * time.Second,
		}

		result, _, err := p.v1api.QueryRange(ctx, query, timeRange)
		if err != nil {
			return nil, err
		}

		data, ok := result.(model.Matrix)
		if !ok {
			return nil, fmt.Errorf("unsupported result format: %s", result.Type().String())
		}

		for _, d := range data {
			var timestamps []int64
			var points []float64
			var labels map[string]string
			labels = make(map[string]string)

			for key, value := range d.Metric {
				labels[string(key)] = string(value)
			}

			for _, value := range d.Values {
				timestamps = append(timestamps, value.Timestamp.Unix())
				points = append(points, float64(value.Value))
			}

			queriesData = append(queriesData, Data{
				Labels:     labels,
				Timestamps: timestamps,
				Points:     points,
			})
		}
	}

	return queriesData, nil
}
