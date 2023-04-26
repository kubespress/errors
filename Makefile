##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: license-header-checker ## Run go fmt against code.
	go fmt ./...
	$(LICENSE_HEADER_CHECKER) -a ./hack/boilerplate.go.txt . go

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ginkgo ## Run tests.
	$(GINKGO) -v --coverprofile cover.out -p ./...

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
LICENSE_HEADER_CHECKER ?= $(LOCALBIN)/license-header-checker
GINKGO ?= $(LOCALBIN)/ginkgo

## Tool Versions
LICENSE_HEADER_CHECKER_VERSION ?= v1.4.0
GINKGO_VERSION ?= v2.9.2

.PHONY: license-header-checker
license-header-checker: $(LICENSE_HEADER_CHECKER) ## Download license-header-checker locally if necessary. 
$(LICENSE_HEADER_CHECKER): | $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/lluissm/license-header-checker/cmd/license-header-checker@$(LICENSE_HEADER_CHECKER_VERSION)

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary. 
$(GINKGO): | $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@$(GINKGO_VERSION)
