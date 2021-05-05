# Build Pinger binary

FROM golang:1.16.0-alpine3.13 AS builder

WORKDIR /go/src/github.com/sdslabs/pinger

COPY . .

ARG vers

RUN apk update && \
  apk add make && \
  apk add bash
RUN make build VERSION=$vers

# Copy binary into actual image

FROM alpine:3.12.0

WORKDIR /go/bin

COPY --from=builder /go/src/github.com/sdslabs/pinger/pinger .

CMD [ "./pinger", "version" ]
