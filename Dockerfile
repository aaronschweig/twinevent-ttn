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

# Command to run
ENTRYPOINT ["/main"]
