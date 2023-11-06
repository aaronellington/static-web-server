## Build the Go package
FROM golang:1.21-bullseye as goBuilder
WORKDIR /workspace/
COPY . .
RUN make clean build-go

FROM debian:bullseye
WORKDIR /workspace/
COPY --from=goBuilder /workspace/var/build ./build
CMD ["./build"]
EXPOSE 2828
