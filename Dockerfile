FROM golang:1.22-alpine as builder
WORKDIR /build

COPY . .
RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o manager cmd/manager/main.go


FROM alpine:3
WORKDIR /workspace

RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=builder /build/manager manager
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./manager"]
