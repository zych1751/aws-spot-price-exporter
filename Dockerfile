FROM golang as builder
COPY . /
RUN GOOS="$(echo "$TARGETPLATFORM" | cut -d "/" -f 1)" \
    GOARCH="$(echo "$TARGETPLATFORM" | cut -d "/" -f 2)" \
    CGO_ENABLED=0 \
    go build -o /app /cmd/main.go

FROM scratch

COPY --from=builder /app /app
USER 1000

ENTRYPOINT ["/app"]
