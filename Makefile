CLI_NAME := tale-gmail

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
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install
install: ## Install
	go get github.com/donofden/tale-gmail/pkg/gmailparser

.PHONY: build
build: ## Build the CLI
	go build
