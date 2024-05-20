FROM golang:1.22 as builder
WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o manager cmd/manager/main.go


FROM alpine:3
WORKDIR /workspace

RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN addgroup --gid 1000 app
RUN adduser --uid 1000 --ingroup app --home /workspace --disabled-password app

COPY --from=builder /build/manager manager

USER app

ENTRYPOINT ["./manager"]
