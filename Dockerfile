## Build the Go package
FROM golang:1.21-buster as goBuilder
WORKDIR /workspace/
COPY . .
RUN make clean build-go

FROM debian:buster
WORKDIR /workspace/
COPY --from=goBuilder /workspace/var/build ./build
CMD ["./build"]
EXPOSE 2828
