# kontrol [![CI](https://github.com/frantjc/kontrol/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/frantjc/kontrol/actions) [![godoc](https://pkg.go.dev/badge/github.com/frantjc/kontrol.svg)](https://pkg.go.dev/github.com/frantjc/kontrol) [![goreportcard](https://goreportcard.com/badge/github.com/frantjc/kontrol)](https://goreportcard.com/report/github.com/frantjc/kontrol) ![license](https://shields.io/github/license/frantjc/kontrol)

Kontrol is a CLI and Kubernetes controller with the goal of making packaging and deploying other controllers easier. It does this by providing a way to package a controller along with the Kubernetes manifests required for it to run.

It doesn't reinvent the wheel; it intends to operate on the controller's image after something like `docker` has built it, it can be used with `kubebuilder` and it outputs manifests that can be applied with `kubectl`.

## install

From a [release](https://github.com/frantjc/kontrol/releases).

Using `brew`:

```sh
brew install frantjc/tap/kontrol
```

From source:

```sh
git clone https://github.com/frantjc/kontrol
cd kontrol
make
```

Using `go`:

```sh
go install github.com/frantjc/kontrol/cmd/kontrol
```

## usage

### Package

Build your controller's image:

```sh
docker build path/to/controller -t your/tag
```

Bundle your controller's manifests with the image:

```sh
kontrol package your/tag \
    --crds path/to/crds.yaml \
    --roles path/to/role.yaml
```

### Deploy

Apply your controller's manifests:

```sh
kontrol deploy your/tag | kubectl apply -f -
```

### Setup Kontrol

Install:

```yml
  - uses: frantjc/kontrol@0.2
```
