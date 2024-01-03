FROM golang:alpine as builder
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . /
RUN GOOS="$(echo "$TARGETPLATFORM" | cut -d "/" -f 1)" \
    GOARCH="$(echo "$TARGETPLATFORM" | cut -d "/" -f 2)" \
    CGO_ENABLED=0 \
    go build -o /app /cmd/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
USER 1000

CMD ["/app"]
ENTRYPOINT ["/app"]
