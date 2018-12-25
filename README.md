# kyml - Kubernetes YAML

A CLI, which helps you to manage your Kubernetes YAML files for your application.

## Background

There are many tool out there to manage Kubernetes manifests, e.g. [ksonnet](https://ksonnet.io/) or [kustomize](https://github.com/kubernetes-sigs/kustomize). They make manifests DRY by introducing other configuration files. I think at least for smaller applications that's too much and I wanted something simpler. So here is `kyml`.

### Goals

- Easily readable and understandable YAML files.

  This means everything is visible in the Kubernetes YAML files themselves. You don't have to jump between multiple files to see what's going on. There is no additional configuration file for yet another tool.

- Stop accidental drift of files between environments.

  In the beginning you decide how your environments (e.g. production and staging) should look like and how they differ. After that it should be hard to only update one environment (e.g. to add another env var) but forget to update the others.

- Support for dynamic values like image tag or branch name.

  Some values are dynamic. E.g. you may deploy every feature branch and want to include the branch name in a label or namespace. A tool has to support this.

### Approach

- Duplicate all files for each environment.
- Lint diff between files with automated tests.
- Simple templating for dynamic values.

## Install

TODO ...

## Usage

```console
$ kyml help
kyml helps you to manage your Kubernetes YAML files.

Usage:
  kyml [command]

Available Commands:
  cat         Concatenate Kubernetes YAML files to stdout
  completion  Generate completion scripts for your shell
  help        Help about any command
  resolve     Resolve image tags to their distribution digest
  test        Run a snapshot test on the diff between Kubernetes YAML files of two environments
  tmpl        Template Kubernetes YAML files

Flags:
  -h, --help      help for kyml
      --version   version for kyml

Use "kyml [command] --help" for more information about a command.
```

### 1. Structure your manifests in the way you want

For most of the following examples, we assume this structure:

```
manifests
├── staging
│   ├── deployment.yaml
│   ├── ingress.yaml
│   └── service.yaml
└── production
    ├── deployment.yaml
    ├── ingress.yaml
    └── service.yaml
```

And some of them use this one:

```
manifests
├── base
│   ├── ingress.yaml
│   └── service.yaml
└── overlays
    ├── staging
    │   └── deployment.yaml
    └── production
        └── deployment.yaml
```

You can adapt these or use anything else, that makes sense for your application.

### 2. Concatenate your files

In the simplest case you concatenate your files and pipe them into `kubectl apply` to deploy them. This does 2 things:

- If multiple files contain the same Kubernetes resource, `kyml cat` deduplicates them. Only the one specified last makes it into the output.
- Resources are sorted by dependencies. So even if you run `kyml cat deployment.yml namespace.yml` the namespace will appear first in the output. This way you don't have to prefix your filenames with numbers just to make sure your resources are created in the correct order.

```sh
kyml cat manifests/production/* | kubectl apply -f -
```

```sh
kyml cat manifests/base/* manifests/overlays/production/* | kubectl apply -f -
```

### 3. Test that your environment don't drift apart

You can add a `kyml test` to this pipeline. This will create a diff between your environments. If this diff does not match your stored snapshot, the command fails and nothing gets deployed.

```sh
kyml cat manifests/production/* |
    kyml test manifests/staging/* \
        --name-main production \
        --name-comparison staging \
        --snapshot-file tests/snapshot-production-vs-staging.diff |
    kubectl apply -f -
```

### 4. Inject dynamic values

Inject dynamic values using `kyml tmpl`, which supports the [go template](https://golang.org/pkg/text/template/) syntax. All values you provide on the command line will be replaced in the pipeline. For example to replace `{{.ImageTag}}` in your manifests, specify the flag `--value ImageTag=my-dynamic-tag`.

```sh
kyml cat manifests/production/* |
    kyml test manifests/staging/* \
        --name-main production \
        --name-comparison staging \
        --snapshot-file tests/snapshot-production-vs-staging.diff |
    kyml tmpl \
        -v Greeting=hello \
        -v ImageTag=$(git rev-parse --short HEAD) \
        -e TRAVIS_BRANCH |
    kubectl apply -f -
```

### 5. Resolve Docker images to their digest

If you tag the same image multiple times (e.g. because you build every commit, tagging images with the commit sha), you may want to resolve the tags to the image digest. This makes sure Kubernetes only restarts your applications if the image content changed.

```sh
kyml cat manifests/production/* |
    kyml tmpl -v ImageTag=$(git rev-parse --short HEAD) |
    kyml resolve |
    kubectl apply -f -
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
