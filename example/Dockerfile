# syntax=docker/dockerfile:1.2

FROM registry.suse.com/bci/golang:1.18 as builder

RUN zypper --non-interactive up

RUN mkdir /app
WORKDIR /app

COPY ../go.mod ./
COPY ../go.sum ./

RUN go mod download

COPY .. .

RUN go build -o openfga-scim-bridge ./cmd/main.go

FROM registry.suse.com/bci/bci-base:latest

RUN zypper --non-interactive up

COPY --from=builder /app/openfga-scim-bridge /usr/local/bin/openfga-scim-bridge

RUN mkdir /app
WORKDIR /app

EXPOSE 8080

CMD ["/usr/local/bin/openfga-scim-bridge", "server"]
