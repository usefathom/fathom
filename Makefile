EXECUTABLE := fathom
LDFLAGS += -extldflags "-static" -X "main.version=$(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')" -X "main.commit=$(shell git rev-parse HEAD)"  -X "main.date=$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
MAIN_PKG := ./main.go
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
ASSET_SOURCES ?= $(shell find assets/src/. -type f)
GO_SOURCES ?= $(shell find . -name "*.go" -type f)
GOPATH=$(shell go env GOPATH)
ARCH := amd64
OS := linux

.PHONY: all
all: build 

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(GO_SOURCES) assets/build
	go build -o $@ $(MAIN_PKG)

.PHONY: docker
docker: $(GO_SOURCES) $(GOPATH)/bin/packr2
	GOOS=$(OS) GOARCH=$(ARCH) $(GOPATH)/bin/packr2 build -v -ldflags '-w $(LDFLAGS)' -o $(EXECUTABLE) $(MAIN_PKG)

$(GOPATH)/bin/packr2:
	GOBIN=$(GOPATH)/bin go install github.com/gobuffalo/packr/v2/packr2@latest

.PHONY: npm 
npm:
	if [ ! -d "node_modules" ]; then npm install; fi

assets/build: $(ASSET_SOURCES) npm
	./node_modules/gulp/bin/gulp.js	

assets/dist: $(ASSET_SOURCES) npm
	NODE_ENV=production ./node_modules/gulp/bin/gulp.js

.PHONY: clean
clean:
	go clean -i ./...
	$(GOPATH)/bin/packr clean
	rm -rf $(EXECUTABLE)

.PHONY: fmt
fmt:
	go fmt $(PACKAGES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: errcheck
errcheck:
	@which errcheck > /dev/null; if [ $$? -ne 0 ]; then \
		go install github.com/kisielk/errcheck@latest; \
	fi
	errcheck $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		go install github.com/golang/lint/golint@latest; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: test
test:
	for PKG in $(PACKAGES); do go test $$PKG || exit 1; done;

.PHONY: referrer-spam-blacklist
referrer-spam-blacklist:
	wget https://raw.githubusercontent.com/matomo-org/referrer-spam-blacklist/master/spammers.txt -O pkg/aggregator/data/blacklist.txt
	go-bindata -prefix "pkg/aggregator/data/" -o pkg/aggregator/bindata.go -pkg aggregator pkg/aggregator/data/
