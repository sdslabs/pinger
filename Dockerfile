# Build Pinger binary

# Use the golang image for building the binary
FROM golang:1.16.0-alpine3.13 AS builder

# Set the work directory
WORKDIR /go/src/github.com/sdslabs/pinger

# Copy over the source code to the container
COPY . .

# Install bash and make
RUN apk update && \
  apk add make && \
  apk add bash

# Build the actual binary
RUN make build VERSION=1.0

# Copy binary into actual image

# Use the alpine image for running the binary
FROM alpine:3.13

# Set the work directory
WORKDIR /go/bin

# Copy over the binary from Builder image
COPY --from=builder /go/src/github.com/sdslabs/pinger/pinger .

# Copy over the agent.yml from Builder image
COPY --from=builder /go/src/github.com/sdslabs/pinger/agent.yml .

# Final command to run pinger
CMD [ "./pinger", "agent" ]
