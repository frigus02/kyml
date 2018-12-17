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
kyml cat feature/* |
    kyml tmpl \
        -v Greeting=hello \
        -v ImageTag=$(git rev-parse --short HEAD) \
        -e TRAVIS_BRANCH |
    kubectl apply -f -

kyml cat base/* overlays/production/* | ...

kyml cat production/*.yml |
    kyml test staging/* \
        --name-main production \
        --name-comparison staging \
        --snapshot-file tests/snapshot-production-vs-staging.diff |
    kyml tmpl \
        -v Greeting=hello \
        -v ImageTag=$(git rev-parse --short HEAD) \
        -e TRAVIS_BRANCH |
    kubectl apply -f -
```

```sh
# Use https://goreleaser.com/ ?!
go install -ldflags "-X github.com/frigus02/kyml/pkg/commands.version=0.0.1"
```
