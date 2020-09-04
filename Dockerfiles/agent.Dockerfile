# Copyright (c) 2020 SDSLabs
# Use of this source code is governed by an MIT license
# details of which can be found in the LICENSE file.

FROM golang:1.14.2-alpine3.11

ENV METRICS_DB_NAME "status-dev"
ENV METRICS_HOST "127.0.0.1"
ENV METRICS_PASSWORD "password"
ENV METRICS_PORT "5432"
ENV METRICS_SSL_MODE "true"
ENV METRICS_USERNAME "username"
ENV PORT "9019"
ENV INTERVAL "2m"

WORKDIR /agent
COPY . .

RUN apk update && apk add make
RUN make build

CMD ./target/pinger agent \
    --metrics-db-name=$METRICS_DB_NAME \
    --metrics-host=$METRICS_HOST \
    --metrics-password=$METRICS_PASSWORD \
    --metrics-port=$METRICS_PORT \
    --metrics-ssl-mode=$METRICS_SSL_MODE \
    --metrics-username=$METRICS_USERNAME \
    --port=$PORT \
    --interval=$INTERVAL
