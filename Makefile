BINARY_NAME = jdy
GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test

all: serve

init:
	$(GOMOD) init $(module)

install:
	$(GOMOD) tidy

serve:
	$(GOCMD) run .

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_NAME) -v ./ ;
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_NAME).mac -v ./ ;
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_NAME).exe -v ./

clean:
	$(GOCMD) clean;
	rm -rf ./build/$(BINARY_NAME)*;

.PHONY: serve build
