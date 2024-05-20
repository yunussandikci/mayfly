FROM golang:1.22 as builder
WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o manager cmd/manager/main.go


FROM alpine:3
WORKDIR /workspace

RUN apk --no-cache add ca-certificates && update-ca-certificates

RUN addgroup --gid 1000 app
RUN adduser --disabled-password --gecos "" --home "/workspace" --ingroup app --no-create-home --uid 1000 app

COPY --from=builder /build/manager manager
USER 1000

ENTRYPOINT ["./manager"]
