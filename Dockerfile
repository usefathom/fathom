FROM golang:latest AS compiler
WORKDIR /go/src/github.com/usefathom/fathom
ADD . /go/src/github.com/usefathom/fathom
RUN go get -u github.com/gobuffalo/packr/... && make docker

FROM alpine:latest
EXPOSE 8080
RUN apk add --update --no-cache bash ca-certificates
WORKDIR /app
COPY --from=compiler /go/src/github.com/usefathom/fathom/fathom .
CMD ["./fathom", "server"]
