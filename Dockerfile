ARG ASSETDIST=node:alpine
ARG GODIST=golang:latest
ARG GOARCH=amd64
ARG RELEASEDIST=alpine:latest

FROM $ASSETDIST AS assetbuilder
RUN apk --no-cache add git python build-base
WORKDIR /app
COPY package*.json ./
COPY gulpfile.js ./
COPY assets/ ./assets/
RUN npm install && NODE_ENV=production npx gulp

# Need to persist the build args across build stages
#    By default they are being consumed by the first build step for some reason
#    See docker docs : https://docs.docker.com/engine/reference/builder/#scope
#    "An ARG instruction goes out of scope at the end of the build stage where it was defined. To use
ARG GODIST
ARG GOARCH
ARG RELEASEDIST

FROM $GODIST AS binarybuilder
WORKDIR /go/src/github.com/usefathom/fathom
COPY . /go/src/github.com/usefathom/fathom
COPY --from=assetbuilder /app/assets/build ./assets/build
RUN go get -u github.com/gobuffalo/packr/... && make ARCH=${GOARCH} docker

# Persist build arg (see above for details)
ARG RELEASEDIST

FROM ${RELEASEDIST}
EXPOSE 8080
RUN apk add --update --no-cache bash ca-certificates sqlite libpq postgresql-client mysql-client
WORKDIR /app
COPY --from=binarybuilder /go/src/github.com/usefathom/fathom/fathom .
CMD ["./fathom", "server"]

