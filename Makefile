SCRIPTS_DIR := ./scripts

all: test

fmt:
	go fmt ./...

fmtcheck:
	bash $(SCRIPTS_DIR)/verify-gofmt.sh

tidycheck:
	bash $(SCRIPTS_DIR)/verify-gomod-tidy.sh

verify: fmtcheck tidycheck

test: fmtcheck
	go test -v ./...
