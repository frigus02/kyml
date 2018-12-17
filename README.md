# kyml - Kubernetes YAML

A way of managing Kubernetes YAML files.

## Background

### Goals

- Easly readable YAML files.

  This means everything is visible in the YAML files and you don't have to jump between multiple files to see what's going on.

- Stop accidental drift of files between environments.

  Production, staging or whatever environments you have should stay as similar as possible. We don't want you to accidentally forget to update one of the environments.

- Support for dynamic values like image tag or branch name.

  Some values are dynamic. E.g. you may deploy every feature branch and want to include the branch name in a label or namespace.

### Approach

- Duplicate files for each environment.
- Lint diff between files with automated tests.
- Provide easy to use tool to edit dynamic things like image tag or branch name in YAML files using CLI.

## Usage

```
~/someApp
├── feature
│   ├── deployment.yaml
│   ├── service.yaml
├── staging
│   ├── deployment.yaml
│   ├── ingress.yaml
│   └── service.yaml
└── production
    ├── deployment.yaml
    ├── ingress.yaml
    └── service.yaml
```

```sh
kyml cat staging/* |
    kyml tmpl \
        -e SECRET \
        -v IMAGE=monopole/hello:$(git rev-parse HEAD) \
        -v BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD) |
    kubectl apply -f -

kyml cat base/* overlays/staging/* | ...

kyml cat staging/*.yml |
    kyml test \
        --name-stdin staging \
        --name-files feature \
        --snapshot-file staging/kyml-snapshot-vs-feature.diff \
        feature/* |
    kyml tmpl \
        -e SECRET \
        -v IMAGE=monopole/hello:$(git rev-parse HEAD) \
        -v BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD) |
    kubectl apply -f -
```

```sh
# Use https://goreleaser.com/ ?!
go install -ldflags "-X github.com/frigus02/kyml/pkg/commands.version=0.0.1"
```
