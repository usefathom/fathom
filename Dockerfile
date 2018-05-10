FROM alpine:latest

EXPOSE 8080
WORKDIR /app
ADD ./fathom /app/fathom
CMD ["./fathom"]

RUN apk add --update bash ca-certificates && rm -rf /var/cache/apk/*
