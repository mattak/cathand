GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
BINARY_NAME=cathand
BINARY_DIR=bin

all: test build

.PHONY: test
test:
	$(GOTEST) -v ./test/

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/cathand/main.go

.PHONY: run
run:
	$(GORUN) ./cmd/cathand/main.go

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -r $(BINARY_DIR)

.PHONY: install
install:
	cd cmd/cathand && go install

