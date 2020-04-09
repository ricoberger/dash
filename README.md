# dash

dash is a terminal dashboard solution inspired by [Grafana](https://grafana.com).

![node_exporter](./examples/assets/node_exporter.png)

## Features

- [x] **Multiple Datasources:** Multiple datasources can be defined via yaml files.
- [x] **Multiple Dashboards:** Dashboards can be defined via yaml files and can be switch during runtime.
- [x] **Time Interval:** The initial time interval can be set via the `--config.interval` command-line flag and can be changed during runtime.
- [x] **Refresh Rate:** The initial refresh rate can be set via the `--config.refresh` command-line flag and can be changed during runtime.
- [x] **Multiple Visualizations:** Choose between singlestat, gauge, donut, sparkline and linechart to visualize your data.
- [ ] **Dynamic Datasources:** Currently the datasource must be provided for a graph. A better solution would be to have a datasource variable to select the datasource for a dashboard.
- [ ] **Other Datasources:** Support for datasources besides Prometheus (e.g. InfluxDB, Elasticsearch, ...)

**Note:** If you want to contribute (adding a missing features or a new one) feel free to create a PR. If you want to share a dashboard please add the `.yaml` file and a screenshot to the examples folder.

## Installation

See [https://github.com/ricoberger/dash/releases](https://github.com/ricoberger/dash/releases) for the latest release.

```sh
VERSION=
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
wget https://github.com/ricoberger/dash/releases/download/$VERSION/dash-$GOOS-$GOARCH
sudo install -m 755 dash-$GOOS-$GOARCH /usr/local/bin/dash
```

## Usage

By default dash will look at `~/.dash` for the datasources and dashboards. To change the location of the datasources and dashboards you can pass the folder via the `--config.dir` command-line flag. Inside the configuration folder dash loads all datasources from the `datasources` and dashboards from the `dashboards` folder.

```
~/.dash
├── dash.log
├── dashboards
│   └── node_exporter.yaml
└── datasources
    └── prometheus.yaml
```

To set the initial interval and refresh rate you can use the `--config.interval` and `--config.refresh` flag. To enable the debug logs you can pass the `--debug` flag to dash.

### Datasources

The configuration file for a datasource looks as follows:

```yaml
# Type is the type of the datasource. Currently only "Prometheus" is supported.
type: <string>
# Custom name of the datasource, which is used in the dashboards.
name: <string>
# URL endpoint for the datasource.
url: <string>
# Credentials for authentication. You can use basic authentication via "username" and "password" or token authentication via the "token" field.
auth:
  username: <string>
  password: <string>
  token: <string>
# Type specific options for a datasource.
options:
  # Maximum number of points which should be returned for a time series from Prometheus.
  maxPoints: <int>
  # Steps between points in a Prometheus time series in seconds. The value is only used if the "maxPoints" option is not set. The default value is 10 seconds.
  step: <int>
```

### Dashboards

You can take a look at the [examples folder](./examples) for some predefined dashboards. The configuration file for a dashboards looks as follows:

```yaml
# Name of the dashboard.
name: <string>
# Variable in the dashboard.
variables: [ <variable> ]
# Rows in the dashboard.
rows: [ <row> ]
```

#### Variable

```yaml
# Name of the variable, which can be used in queries.
name: <string>
# Name of the datasource, must be defined in a datasource file.
datasource: <string>
# Query which is executed against the datasource.
query: <string>
# Label which should be used to for the values of the variable.
label: <string>
# Add an include all option to the list of variable values.
all: <bool>
```

#### Row

```yaml
# Height of the row. Must be a value between 0 and 1. The sum of the height must be 1.
height: <int>
# Graphs which should be shown in the row.
graphs: <graph>
```

#### Graph

```yaml
# Width of the graph in the row.
width: <int>
# Name of the datasource, were the queries should run against.
datasource: <string>
# Visualization type. Must be "singlestat", "gauge", "donut", "sparkline" or "linechart".
type: <string>
# Title of the visualization.
title: <string>
# List of queries which should be used to retrive the data for the visualization.
queries: [ <query> ]
# Visualization type specific configuration.
options:
  # Unit which should be show for the metric.
  # This option is available for "singlestat", "sparkline" and "linechart".
  unit: <string>
  # Array of stats which should be shown. Possible values are "current", "first", "min", "max", "avg", "total", "diff" and "range".
  # This option is available for "singlestat", "sparkline" and "linechart".
  # "singlestat": If this is not provided, "current" will be used. If the length of the array is greater then 1 only the first value will be used.
  # "sparkline": The "current" value will always be shown. You can add multiple other values which should be shown in the legend.
  # "linechart": The "current" value will always be shown. You can add multiple other values which should be shown in the legend.
  stats: [ <string> ]
  # Decimals for the metric. The default value is 0.
  # This option is available for "singlestat", "sparkline" and "linechart".
  decimals: <int>
  # Threshold which are used for the different colors. The length of the array must be the length of colors - 1.
  # This option is available for "singlestat", "gauge", "donut" and "sparkline".
  # By default the value is compared with the "current" value. If you have set a value for "stats" it is compared against the first value in the provided array.
  thresholds: [ <float> ]
  # Colors for the metrics. The array length must be the length of thresholds + 1.
  colors: [ <string> ]
  # Mappings allow you to overwrite the returned value for a "singlestat"
  mappings <map[string]string>
```

#### Query

```yaml
# Query which should be executed against the configured datasource.
query: <string>
# Label which should be used for the query. This must be unique per returned time series. Returned labels can be used via templating, e.g. "trans {{.device}}".
label: <string>
```

### Key Mapping

Use the following keys to navigate within dash:

| Key | Function |
| --- | -------- |
| `q`, `<Ctrl-C>` | Close dash. |
| `<Esc>` | Close modal with out changing the value. |
| `<Enter>` | Select the current value in the modal and close the modal. |
| `<0-9>` | Provide the index for a value in the modal. |
| `d` | Open the dashboards modal. |
| `v` + `<1-9>` | Open the modal for a variable. Press `v` plus a number from `1` to `9` to select the variable. |
| `i` | Open the interval modal. |
| `r` | Open the refresh rate modal. |
