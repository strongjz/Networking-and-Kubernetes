FROM golang:1.15 AS builder
WORKDIR /opt
COPY go.mod .
COPY web-server.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o web-server .

FROM golang:1.15
WORKDIR /opt
COPY --from=0 /opt/web-server .
CMD ["/opt/web-server"]
