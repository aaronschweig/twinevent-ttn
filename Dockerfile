FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

# Build a small image
FROM alpine

COPY --from=builder /dist/main /
COPY --from=builder /build/ditto/*.json /ditto/
COPY --from=builder /build/default-config.yaml /default-config.yaml

# Command to run
ENTRYPOINT ["/main"]

CMD ["-config", "./default-config.yaml"]
