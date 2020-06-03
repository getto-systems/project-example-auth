FROM golang:1.14.4-buster as builder
COPY . /build
WORKDIR /build
RUN : && \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build -a -o app . && \
  :

FROM scratch
COPY --from=builder /build/app /app
CMD ["/app"]
