FROM golang:1.12-alpine as builder

RUN apk add --no-cache git

ADD . /gopath/src/github.com/astranet/meshRPC-benchmark
ENV GOPATH=/gopath
RUN go get github.com/astranet/meshRPC-benchmark/bench_api

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /gopath/bin/bench_api /usr/local/bin/

EXPOSE 11999
EXPOSE 8282
ENTRYPOINT ["bench_api"]
