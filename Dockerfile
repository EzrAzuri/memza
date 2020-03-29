FROM golang:1.13 as build
WORKDIR /build
COPY . .
RUN cd ./cmd/memza-api && GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o memza ./main.go

FROM alpine:latest
LABEL maintainer="Ron Compos <rcompos@gmail.com>"
COPY --from=build /build/cmd/memza-api/memza /bin
EXPOSE 8080
ENV MEMCACHED_SERVER_URL 0.0.0.0:11211
ENTRYPOINT ["memza"]
CMD ["-d"]
