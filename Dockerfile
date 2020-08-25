FROM golang:1.15.0-buster as builder
COPY . /build
WORKDIR /build
RUN : && \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build -a -o app $(head -1 go.mod | cut -d' ' -f2)/x_http_server && \
  :

FROM gcr.io/distroless/static-debian10
COPY --from=builder /build/app /app
CMD ["/app"]
