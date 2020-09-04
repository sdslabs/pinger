# Copyright (c) 2020 SDSLabs
# Use of this source code is governed by an MIT license
# details of which can be found in the LICENSE file.

FROM golang:1.15.1-alpine3.12

ENV METRICS_DB_NAME "<db_name>"
ENV METRICS_HOST "<db_host>"
ENV METRICS_PORT "<db_port>"
ENV METRICS_USERNAME "<db_username>"
ENV METRICS_PASSWORD "<db_password>"
ENV METRICS_SSL_MODE "<db_in_ssl_mode>"
ENV PORT "<agent_port>"
ENV INTERVAL "<update_interval>"

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
