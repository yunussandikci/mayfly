FROM golang:1.22 as builder
WORKDIR /build

COPY . .
RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o manager cmd/manager/main.go

FROM scratch
WORKDIR /workspace

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/manager manager

ENTRYPOINT ["./manager"]
