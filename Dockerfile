FROM golang:latest as build
WORKDIR /streambot
# Copy source code
COPY . .

# Build statically linked binary
RUN go build

###
FROM golang:latest

WORKDIR /

COPY --from=build /streambot/streambot .
ENTRYPOINT ["/streambot"]