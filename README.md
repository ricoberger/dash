# dash

dash is a terminal dashboard solution inspired by [Grafana](https://grafana.com).

![node_exporter](./examples/assets/node_exporter.png)

## Features

- [x] **Multiple Datasources:** Multiple datasources can be defined via yaml files.
- [x] **Multiple Dashboards:** Dashboards can be defined via yaml files and can be switch during runtime.
- [x] **Time Interval:** The initial time interval can be set via the `--config.interval` command-line flag and can be changed during runtime.
- [x] **Refresh Rate:** The initial refresh rate can be set via the `--config.refresh` command-line flag and can be changed during runtime.
- [x] **Multiple Visualizations:** Choose between singlestat, gauge, sparkline and plots to visualize your data.
- [ ] **Show time at X-Axis:** Currently the default values from termui are shown at the x-axis. This should be changed to the time values.
- [ ] **Missing Values:** The handling of missing values in the timeseries can be improved.
- [ ] **Other Datasources:** Support for datasources besides Prometheus (e.g. InfluxDB, Elasticsearch, ...)
- [ ] **Labeling:** Currently only labels from the returned data can be used in the legend. There also custom text should be supported.

**Note:** If you want to contribute, to add one of the missing features or to add a new feature, feel free to create a PR. If you want to share a dashboard please add the `.yaml` file and a screenshot to the examples folder.

## Installation

See [https://github.com/ricoberger/dash/releases](https://github.com/ricoberger/dash/releases) for the latest release.

```sh
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
wget https://github.com/ricoberger/dash/releases/download/0.9.0/dash-$GOOS-$GOARCH
sudo install -m 755 dash-$GOOS-$GOARCH /usr/local/bin/dash
```

## Usage

By default dash will look at `~/.dash` for the datasources and dashboards. To change the location of the datasources and dashboards you can pass the folder via the `--config.dir` command-line flag. Inside the configuration folder dash loads all datasources from the `datasources` and dashboards from the `dashboards` folder.

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
```

### Dashboards

For example dashboards you can take a look at the [examples folder](./examples). The configuration file for a dashboards looks as follows:

```yaml
# Name of the dashboard.
name: <string>
# Variable in the dashboard.
variables: [ <variable> ]
# Rows in the dashboard.
rows: [ <row> ]
```

```yaml
# Name of the variable, which can be used in queries.
name: <string>
# Name of the datasource, must be defined in a datasource file.
datasource: <string>
# Query which is executed against the datasource.
query: <string>
# Label which should be used to for the values of the variable.
label: <string>
```

```yaml
# Height of the row. Must be a value between 0 and 1. The sum of the height must be 1.
height: <float>
# Graphs which should be shown in the row.
graphs: <graph>
```

```yaml
# Width of the graph in the row.
width: <float>
# Name of the datasource, were the queries should run against.
datasource: <string>
# Visualization type. Must be "singlestat", "gauge", "sparkline" or "plot"
type: <string>
# Title of the visualization.
title: <string>
# Queries which should be used in the visualization.
queries: [ <string> ]
options:
  # Unit which should be show for the metric.
  unit: <string>
  # Array of stats which should be shown. Possible values are "current", "first", "min", "max", "avg", "total", "diff" and "range".
  stats: [ <string> ]
  # Prefix is the text which is shown before the value in a singlestat.
  prefix: <string>
  # Postfix is the text which is shown behind the value in a singlestat.
  postfix: <string>
  # Decimals for the metric. The default value is 0.
  decimals: <int>
  # Threshold which are used for the different colors. The length of the array must be the length of colors - 1.
  thresholds: [ <float> ]
  # Colors for the metrics. The array length must be the length of thresholds + 1.
  colors: [ <string> ]
  # Label from the query result which should be used for the legend.
  label: <string>
```

### Key Mapping

Use the following keys to navigate within dash:

| Key | Function |
| --- | -------- |
| `q`, `<Ctrl-C>` | Close dash. |
| `k`, `<Up>` | Scroll up in the open modal. |
| `j`, `<Down>` | Scroll down in the open modal. |
| `<Esc>` | Close modal with out changing the value. |
| `<Enter>` | Select the current value in the modal and close the modal. |
| `d` | Open the dashboards modal. |
| `v` + `<1-9>` | Open the modal for a variable. Press `v` plus a number from `1` to `9` to select the variable. |
| `i` | Open the interval modal. |
| `r` | Open the refresh rate modal. |
