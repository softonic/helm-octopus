# Helm Octopus Plugin
[![Go Report Card](https://goreportcard.com/badge/github.com/databus23/helm-diff)](https://goreportcard.com/report/github.com/softonic/helm-octopus)
[![GoDoc](https://godoc.org/github.com/databus23/helm-diff?status.svg)](https://godoc.org/github.com/softonic/helm-octopus)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/softonic/helm-octopus/blob/master/LICENSE)

This Helm plugin allows to reference packaged value files (other than the default values.yaml).


### Install
```bash
helm plugin install https://github.com/softonic/helm-octopus
```

### Supported helm commands

Octopus will kick-in only with the following commands:
* upgrade
* install
* template
* lint

It can "handle" all other commands as well, but they will be proxied to helm binary.

### Usage

Octopus will read package value files, as long they have the following format:

`subchart://<subchartId>/<filename>`

Where `subchartId` is dependency chart `alias` or `name`, and filename is the file path
within its package.


#### Example

For example, this is our dependency tree structure:

```bash
mydep
├── charts
├── Chart.yaml
├── templates
│   ├── deployment.yaml
│   └── _helpers.tpl
├── values.yaml
└── values.custom.yaml
```

If we have the following chart:

```yaml
apiVersion: v2
name: mychart
version: 0.0.1
dependencies:
- name: mydep
  alias: foobar
  version: 1.0.0
  repository: "@myrepo"
```

We could use octopus to reference `values.custom.yaml` for `mydep` chart.

`helm octopus template myrelease ./mychart -f subchart://foobar/values.custom.yaml`

#### Env vars

`HELM_OCTOPUS_TMP_DIR` can be used to define temporary directory where to save the package value files.
Defaults to `/tmp/octopus/`.
