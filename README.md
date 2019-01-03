# kyml - Kubernetes YAML

A CLI, which helps you to work with and deploy plain Kubernetes YAML files.

## Background

There are many great tools out there to manage Kubernetes manifests, e.g. [ksonnet](https://ksonnet.io/) or [kustomize](https://github.com/kubernetes-sigs/kustomize). They try to make working with manifests easier by deduplicating config. However they usually introduce other configuration files, which comes with complexity on its own. I wanted something simpler, especially for smaller applications.

So here is `kyml`:

- Work with plain Kubernetes YAML files. No additional config files.
- Duplicate files for each environment. But ensure updates always happen to all environments.
- Support dynamic values with limited templating.
- Save to run. Never touch original YAML files.

## Install

Download a binary from the [release page](https://github.com/frigus02/kyml/releases).

This downloads the latest version for Linux:

```sh
curl -sfL -o /usr/local/bin/kyml https://github.com/frigus02/kyml/releases/download/v20190103/kyml_20190103_linux_amd64 && chmod +x /usr/local/bin/kyml
```

## Usage

`kyml` provides commands for concatenating YAML files, testing and templating. Commands usually read manifests from stdin and print them to stdout. This allows you to build a pipeline of commands, which can end with a `kubectl apply` to do the deployment.

- [Structure your manifests in the way you want](#structure-your-manifests-in-the-way-you-want)
- [`kyml cat` - concatenate YAML files](#kyml-cat---concatenate-yaml-files)
- [`kyml test` - ensure updates always happen to all environments](#kyml-test---ensure-updates-always-happen-to-all-environments)
- [`kyml tmpl` - inject dynamic values](#kyml-tmpl---inject-dynamic-values)
- [`kyml resolve` - resolve Docker images to their digest](#kyml-resolve---resolve-docker-images-to-their-digest)

Run `kyml --help` for details about the different commands.

### Structure your manifests in the way you want

For most of the examples in this readme we assume the following structure:

```
manifests
|- staging
|  |- deployment.yaml
|  |- ingress.yaml
|  `- service.yaml
`- production
   |- deployment.yaml
   |- ingress.yaml
   `- service.yaml
```

And some of them use this:

```
manifests
|- base
|  |- ingress.yaml
|  `- service.yaml
`- overlays
   |- staging
   |  `- deployment.yaml
   `- production
      `- deployment.yaml
```

You can adapt these or use anything else, that makes sense for your application.

### `kyml cat` - concatenate YAML files

Concatenate your files and pipe them into [`kubectl apply`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#apply) to deploy them. This does 2 things:

- If multiple files contain the same Kubernetes resource, `kyml cat` deduplicates them. Only the one specified last makes it into the output.
- Resources are sorted by dependencies. So even if you specify the namespace last (e.g. `kyml cat deployment.yaml namespace.yaml`) the namespace will appear first in the output. This makes sure your resources are created in the correct order.

```sh
kyml cat manifests/production/* | kubectl apply -f -
```

```sh
kyml cat manifests/base/* manifests/overlays/production/* | kubectl apply -f -
```

### `kyml test` - ensure updates always happen to all environments

Testing works by creating a diff between two environments and storing it in a snapshot file. The command compares the diff result to the snapshot and fails if it doesn't match.

`kyml test` reads manifests of the main environment from stdin and files from the comparison environment are specified as arguments, similar to `kyml cat`. If the snapshot matches, it prints the main environment manifests to stdout. This way you can include a test in your deployment command pipeline to make sure nothing gets deployed if the test fails.

```sh
kyml cat manifests/production/* |
    kyml test manifests/staging/* \
        --name-main production \
        --name-comparison staging \
        --snapshot-file tests/snapshot-production-vs-staging.diff |
    kubectl apply -f -
```

### `kyml tmpl` - inject dynamic values

Use templates (in the [go template](https://golang.org/pkg/text/template/) syntax) to inject dynamic values. To make sure values are escaped properly and this feature doesn't get misused you can only template string scalars. Example:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: the-namespace
  labels:
    branch: "{{.TRAVIS_BRANCH}}"
```

`kyml test` reads manifests from stdin and prints the result to stdout. Values are provided as command line options. Use `--value key=value` for literal strings and `--env ENV_VAR` for environment variables. These options can be repeated multiple times. The command fails if the manifests contain any template key, which is not specified on the command line.

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

### `kyml resolve` - resolve Docker images to their digest

If you tag the same image multiple times (e.g. because you build every commit and tag images with the commit sha), you may want to resolve the tags to the image digest. This way Kubernetes only restarts your applications if the image content has changed.

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
