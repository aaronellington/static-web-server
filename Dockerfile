## Build the Go application
FROM golang:1.21-bullseye as goBuilder
WORKDIR /workspace/
COPY . .
RUN make clean
RUN make build-go

## Copy the Go binary to it's own image to run
FROM debian:bullseye
WORKDIR /workspace/
COPY --from=goBuilder /workspace/var/build ./build
CMD ["./build"]
EXPOSE 2828
