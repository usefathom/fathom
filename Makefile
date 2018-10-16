DIST := build
EXECUTABLE := fathom
LDFLAGS += -extldflags "-static" -X "main.Version=$(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')"
MAIN_PKG := ./cmd/fathom
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
JS_SOURCES ?= $(shell find assets/src/. -name "*.js" -type f)
GO_SOURCES ?= $(shell find . -name "*.go" -type f)
SQL_SOURCES ?= $(shell find . -name "*.sql" -type f)
ENV ?= $(shell export $(cat .env | xargs))
GOPATH=$(shell go env GOPATH)

$(EXECUTABLE): $(GO_SOURCES) assets/build
	go build -o $@ $(MAIN_PKG)

.PHONY: all
all: build 

.PHONY: install
install: $(wildcard *.go) $(GOPATH)/bin/packr
	$(GOPATH)/bin/packr install -v -ldflags '-w $(LDFLAGS)' $(MAIN_PKG)

.PHONY: build
build: $(EXECUTABLE)

.PHONY: docker
docker: $(GO_SOURCES) 
	GOOS=linux GOARCH=amd64 $(GOPATH)/bin/packr build -v -ldflags '-w $(LDFLAGS)' -o $(EXECUTABLE) $(MAIN_PKG)

.PHONY: dist
dist: assets/dist 
	GOOS=linux GOARCH=amd64 $(GOPATH)/bin/packr build -v -ldflags '-w $(LDFLAGS)' -o build/fathom-linux-amd64 $(MAIN_PKG)
	GOOS=linux GOARCH=arm64 $(GOPATH)/bin/packr build -v -ldflags '-w $(LDFLAGS)' -o build/fathom-linux-arm64 $(MAIN_PKG)	
	GOOS=linux GOARCH=386 $(GOPATH)/bin/packr build -v -ldflags '-w $(LDFLAGS)' -o build/fathom-linux-386 $(MAIN_PKG)

$(GOPATH)/bin/packr:
	GOBIN=$(GOPATH)/bin go get github.com/gobuffalo/packr/...

.PHONY: npm 
npm:
	if [ ! -d "node_modules" ]; then npm install; fi

assets/build: $(JS_SOURCES) npm
	./node_modules/gulp/bin/gulp.js	

assets/dist: $(JS_SOURCES) npm
	NODE_ENV=production ./node_modules/gulp/bin/gulp.js

.PHONY: clean
clean:
	go clean -i ./...
	packr clean
	rm -rf $(EXECUTABLE) $(DIST) 

.PHONY: fmt
fmt:
	go fmt $(PACKAGES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: errcheck
errcheck:
	@which errcheck > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: test
test:
	for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

