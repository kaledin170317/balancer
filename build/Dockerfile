FROM golang:1.24


WORKDIR /go/balancer
COPY ../ ./

RUN go build -o /balancer ./cmd/balancer/main.go \
    && go clean -cache -modcache

EXPOSE 8080


CMD ["/balancer"]