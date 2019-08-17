PROJECTNAME := $(shell basename "$(PWD)")
GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid

.DEFAULT_GOAL := explain
.PHONY: explain
explain:
	### Welcome
	#
	#	 _        _                                  _ _ 
	#	| |_ __ _| | ___        __ _ _ __ ___   __ _(_) |
	#	| __/ _` | |/ _ \_____ / _` | '_ ` _ \ / _` | | |
	#	| || (_| | |  __/_____| (_| | | | | | | (_| | | |
	# 	 \__\__,_|_|\___|      \__, |_| |_| |_|\__,_|_|_|
	#	                       |___/                     
	#
	#
	### Installation
	#
	# $$ make install
	#
	### Targets
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the CLI in Project Folder
	go build

.PHONY: build-project
build-project: ## Build the CLI in PROJECT bin folder
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

.PHONY: get
get: ## Get Go Dependencies
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get .

.PHONY: install
install: ## Install "go install GOFILES"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

.PHONY: run
run: ## Run the Project
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

.PHONY: restart ## Restart the Application
restart: stop clean build start

.PHONY: start
start: ## Start the application "bin/PROJECT"
	@echo "Starting bin/$(GONAME)"
	@./bin/$(GONAME) & echo $$! > $(PID)

.PHONY: stop
stop: ## Stop the application  "bin/PROJECT"
	@echo "Stopping bin/$(GONAME)"
	@-kill `cat $(PID)` || true

.PHONY: clean
clean: ## Clean the Project
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
