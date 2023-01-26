# Copyright 2020 The Ledger Authors
#
# Licensed under the AGPL, Version 3.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.gnu.org/licenses/agpl-3.0.en.html
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJ=tokenapi-go
ORG_PATH=github.com/danielnegri
REPO_PATH=$(ORG_PATH)/$(PROJ)

DOCKER_IMAGE=$(PROJ)

$( shell mkdir -p bin )
$( shell mkdir -p release/bin )
$( shell mkdir -p results )

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)/bin
# Prefer ./bin instead of system packages for things like protoc, where we want
# to use the version the application uses, not whatever a developer has installed.
export PATH=$(GOBIN):$(shell printenv PATH)

# Version
VERSION ?= $(shell ./scripts/git-version)
COMMIT_HASH ?= $(shell git rev-parse HEAD 2>/dev/null)
BUILD_TIME ?= $(shell date +%FT%T%z)

LD_FLAGS="-s -w -X $(REPO_PATH)/version.CommitHash=$(COMMIT_HASH) -X $(REPO_PATH)/version.Version=$(VERSION)"

# Inject .env file
-include .env
export $(shell sed 's/=.*//' .env)

build: clean bin/ledger

bin/ledger:
	@echo "Building Ledger: ${COMMIT_HASH}"
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/ledger

clean:
	@echo "Cleaning binary folders"
	@rm -rf bin/*
	@rm -rf release/*
	@rm -rf results/*

release-binary:
	@echo "Releasing binary files: ${COMMIT_HASH}"
	@go build -race -o release/bin/ledger -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/ledger

start: build
	@echo "Starting Ledger server"
	@bin/ledger serve

.PHONY: docker-image
docker-image: clean
	@echo "Building $(DOCKER_IMAGE) image"
	@docker build -t $(DOCKER_IMAGE) --rm -f Dockerfile .

test:
	@echo "Testing"
	@go test -v --short ./...

testcoverage:
	@echo "Testing with coverage"
	@mkdir -p results
	@go test -v $(REPO_PATH)/... | go2xunit -output results/tests.xml
	@gocov test $(REPO_PATH)/... | gocov-xml > results/cobertura-coverage.xml

testrace:
	@echo "Testing with Race Detection"
	@go test -v --race $(REPO_PATH)/...

vet:
	@echo "Running go tool vet on packages"
	go vet $(REPO_PATH)/...

fmt:
	@echo "Running gofmt on package sources"
	go fmt $(REPO_PATH)/...

testall: testrace vet fmt testcoverage

.PHONY: fmt \
		release-binary \
		test \
		testall \
		testcoverage \
		testrace \
		vet
