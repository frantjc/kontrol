GO = go
GIT = git
GOLANGCI-LINT = golangci-lint
GORELEASER = goreleaser
INSTALL = sudo install
DOCKER = docker
DOCKER-COMPOSE = $(DOCKER) compose
CONTROLLER-GEN = controller-gen
KUBECTL = kubectl
HELM = helm

BIN ?= /usr/local/bin

SEMVER ?= 0.1.0

manifests:
	@$(CONTROLLER-GEN) rbac:roleName=kontroller crd webhook paths="./..." output:dir=manifests
	@$(CONTROLLER-GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

fmt generate test:
	@$(GO) $@ ./...

download tidy vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

release:
	@$(GIT) tag v$(SEMVER)
	@$(GIT) push --tags

gen: generate
dl: download
ven: vendor
ver: verify
format: fmt

.PHONY: build dl download fmt format gen generate lint manifests release test ven vendor ver verify
