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

SEMVER ?= 0.2.1

.DEFAULT: install

install: build
	@$(INSTALL) ./dist/kontrol_$(GOOS)_$(GOARCH)*/kontrol $(BIN)

build:
	@$(GORELEASER) release --snapshot --clean --skip-docker

.github/action:
	@cd .github/action && $(YARN) all

manifests:
	@$(CONTROLLER-GEN) rbac:roleName=kontroller crd webhook paths="./..." output:dir=$@
	@$(CONTROLLER-GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

generate test:
	@$(GO) $@ ./...

fmt:
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
