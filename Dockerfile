FROM alpine:latest

EXPOSE 8080

WORKDIR /app
VOLUME ["/app/storage"]
CMD ["/app/ana"]

RUN apk add --update bash ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app/storage/sessions && chmod 777 /app/storage
ADD ./static /app/static
ADD ./views /app/views
ADD ./ana /app/ana
