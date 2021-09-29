go-cat
======
[![Go Reference](https://pkg.go.dev/badge/gitlab.com/sorcero/community/go-cat.svg)](https://pkg.go.dev/gitlab.com/sorcero/community/go-cat)

`go-cat` is a Go command line tool which helps keep track
of infrastructure across multiple clouds in a single git repository.
It helps manage multiple components with their API endpoints, with the 
time they were last deployed on, and the commit SHA. This makes tracking
huge infrastructure trees easier on the long run.

The project is in early alpha and under active development.

Installation ✨
---------------

```bash
cd cmd/go-cat
go build .
./go-cat --help
```

Usage 🤔
--------

The `go-cat upsert` command mutates the json file and regenerates a new 
`README.md` with a markdown table containing information. API consumers
can make use of `infra.json` to dynamically check from command line about the 
information about the latest deployment.

The `go-cat add` command adds infrastructure to queue, which can be then pushed to git 
repository using `go-cat push`. This feature is useful for batch updates to
infrastructure, in one transaction.

Workflow 🔧
----------
A general implementation workflow in an infrastructure repository
would be 
```bash

# do infrastructure changes 
terraform apply

terraform output -format=json > infra_info.json
 
go-cat upsert \
  --name company1-api1 --commit-sha=1c32fb0 \
  --cloud=GCP --cloud-project-id=my-awesome-project-41313 \
  --type=run.googleapis.com \
  --subsystem=company1 \
  --deployment-link=https://example.com \
  --git.url=https://gitlab.com/something/xxx.git \
  --git.username=username \
  --git.password=mysupersecretpassword

```

The `git.url`, `git.username` and `git.password` can also be configured using
environment variables `GO_CAT_GIT_URL`, `GO_CAT_GIT_USERNAME`, `GO_CAT_GIT_PASSWORD`
respectively.

To remove any component from the infrastructure catalog, `go-cat` now
supports a `remove` subcommand, which helps us remove components using all
the power of regex.

```bash
# to remove only one component
go-cat remove --id 'gcp/project1/subsystem/component'
# this will remove component named 'component', belonging to 'subsystem', in project 'project1' deployed on GCP.

# to remove all components, subsystems, projects of a cloud
go-cat remove --id 'gcp/*'

# to remove all components which are called 'api' across all clouds
go-cat remove --id '.*/api'

# to remove a specific subsystem across all cloud projects
go-cat remove --id 'gcp/.*/subsystem/component'
```

Contributing 🔍
---------------
Make sure you adhere to Go formatting guidelines when
contributing to this repository

```bash
go fmt
```

License ⚖️
-----------
This software is licensed under the [Mozilla Public License v2.0](./LICENSE)
