FROM golang:1.16-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .


# Build the Go app
RUN go build -o ./api .

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates
RUN ls
COPY --from=build_base /app/api /app/api
COPY config /app/config
WORKDIR /app

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./api"]