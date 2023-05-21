GO = go
GIT = git
GOLANGCI-LINT = golangci-lint
GORELEASER = goreleaser
CONTROLLER-GEN = controller-gen
YARN = yarn
INSTALL ?= sudo install

BIN ?= /usr/local/bin

GOOS = $(shell $(GO) env GOOS)
GOARCH = $(shell $(GO) env GOARCH)

SEMVER ?= 0.1.4

.DEFAULT: install

install: build
	@$(INSTALL) ./dist/forge_$(GOOS)_$(GOARCH)*/forge $(BIN)

build:
	@$(GORELEASER) release --snapshot --clean

.github/action:
	@cd .github/action && $(YARN) all

manifests:
	@$(CONTROLLER-GEN) rbac:roleName=kontroller crd webhook paths="./..." output:dir=manifests
	@$(CONTROLLER-GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

generate:
	@$(GO) $@ ./...

fmt test:
	@$(GO) $@ ./...
	@cd .github/action && $(YARN) $@

download:
	@$(GO) mod $@
	@cd .github/action && $(YARN)

tidy vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

MAJOR = $(word 1,$(subst ., ,$(SEMVER)))
MINOR = $(word 2,$(subst ., ,$(SEMVER)))

release:
	@cd .github/action && \
		$(YARN) version --new-version $(SEMVER)
	@$(GIT) push
	@$(GIT) tag -f v$(MAJOR)
	@$(GIT) tag -f v$(MAJOR).$(MINOR)
	@$(GIT) push --tags -f

action: .github/action
gen: generate
dl: download
ven: vendor
ver: verify
format: fmt

.PHONY: .github/action action build dl download fmt format gen generate lint manifests release test ven vendor ver verify
