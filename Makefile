DIST := build
EXECUTABLE := fathom
LDFLAGS += -extldflags "-static"
MAIN_PKG := ./cmd/fathom
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
JS_SOURCES ?= $(shell find assets/. -name "*.js" -type f)
SOURCES ?= $(shell find . -name "*.go" -type f)
ENV ?= $(shell export $(cat .env | xargs))

.PHONY: all
all: build 

.PHONY: install
install: $(wildcard *.go)
	packr install -v -ldflags '-w $(LDFLAGS)' $(MAIN_PKG)

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES) $(JS_SOURCES)
	if [ ! -d "node_modules" ]; then npm install; fi
	NODE_ENV=production gulp
	packr build -v -ldflags '-w $(LDFLAGS)' -o $@ $(MAIN_PKG) 

.PHONY: docker
docker:
	docker build -t metalmatze/ana:latest .

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

