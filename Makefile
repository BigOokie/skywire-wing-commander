.PHONY: run run-help check cover
.PHONY: install-linters format

# Compilation output
#BUILD_DIR = build
BIN_DIR = bin
#DOC_DIR = docs
#INCLUDE_DIR = include

#run:  ## Run the skycoin node. To add arguments, do 'make ARGS="--foo" run'.
#	./run.sh ${ARGS}

test: ## Run tests for Wing Commander
	go test ./... -timeout=5m

lint: ## Run linters. Use make install-linters first.
	vendorcheck ./...
	$GOPATH/bin/golangci-lint run --no-config --deadline=3m --disable-all --tests \
		-E golint \
		-E goimports \
		-E varcheck \
		-E unparam \
		-E deadcode \
		-E structcheck \
		-E errcheck \
		-E gosimple \
		-E staticcheck \
		-E unused \
		-E ineffassign \
		-E typecheck \
		-E gas \
		-E megacheck \
		-E misspell \
		./...
	# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	#go vet -all ./...

check: lint test  ## Run tests and linters

cover: ## Runs tests on ./cmd/ with HTML code coverage
	go test -cover -coverprofile=cover.out -coverpkg=github.com/BigOokie/skywire-wing-commander/... ./...
	go tool cover -html=cover.out

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	#go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	# Using v1.10.2
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $GOPATH/bin v1.10.2

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/BigOokie/skywire-wing-commander ./cmd
	goimports -w -local github.com/BigOokie/skywire-wing-commander ./internal

#release-bin: ## Build standalone apps. Use osarch=${osarch} to specify the platform. Example: 'make release-bin osarch=darwin/amd64' Supported architectures are the same as 'release' command.
#	cd $(ELECTRON_DIR) && ./build-standalone-release.sh ${osarch}
#	@echo release files are in the folder of electron/release
