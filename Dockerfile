FROM golang:latest AS compiler
WORKDIR /go/src/github.com/usefathom/fathom
ADD . /go/src/github.com/usefathom/fathom
RUN go get -u github.com/gobuffalo/packr/... && make build

FROM alpine:latest
EXPOSE 8080
WORKDIR /app
COPY --from=compiler /go/src/github.com/usefathom/fathom/fathom .
CMD ["./fathom", "server"]
RUN apk add --update bash ca-certificates && rm -rf /var/cache/apk/*
