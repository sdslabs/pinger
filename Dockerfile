# Copyright (c) 2020 SDSLabs
# Use of this source code is governed by an MIT license
# details of which can be found in the LICENSE file.

#######################
# Build Pinger binary #
#######################

FROM golang:1.15.1-alpine3.12 AS builder

WORKDIR /go/src/github.com/sdslabs/pinger

COPY . .

RUN apk update && apk add make
RUN make build

#################################
# Copy binary into actual image #
#################################

FROM alpine:3.12.0

WORKDIR /go/bin

COPY --from=builder /go/src/github.com/sdslabs/pinger/target/pinger .

CMD [ "./pinger", "ping" ]
