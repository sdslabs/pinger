# Build Pinger binary

FROM golang:1.15.1-alpine3.12 AS builder

WORKDIR /go/src/github.com/sdslabs/pinger

COPY . .

RUN apk update && \
  apk add make && \
  apk add bash
RUN make bin

# Copy binary into actual image

FROM alpine:3.12.0

WORKDIR /go/bin

COPY --from=builder /go/src/github.com/sdslabs/pinger/pinger .

CMD [ "./pinger", "version" ]
