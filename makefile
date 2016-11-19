.PHONY: debug
debug: tracker.go models.go
	go build
	$(GOPATH)/bin/ana
