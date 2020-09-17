---
title: Agent
parent: Architecture
nav_order: 4

---

# Agent

An agent actually pings targets to check if their status is OK or not.
The way an agent works is:

1. It receives the check from the central server through the gRPC API.
1. Each agent has a manager. Manager adds a controller which runs the check
   at regular intervals of time.
1. Another controller is created with all the agents which is responsible
   for updating the database with metrics at regular intervals.
1. At regular intervals of time, agent exports metrics to the database.
   Along with the metrics, agent is also responsible to send alerts.

Elements of an agent are pluggable. An agent can support multiple types of
checkers (TCP, DNS, HTTP etc.), various databases (TimescaleDB, Prometheus)
and alerts (Mail, Slack, Discord etc.)

## Standalone

An agent can also work in standalone mode. This means, an agent can be
configured with a config file and will run independently without any
requirement of it being registered with the app database.

## Metrics

Metrics are stored in a different database than the application database.
This has some benefits:

- Since metrics are very often updated, this database can be scaled.
- Metrics databases are pluggable. Being independent of the main database,
  user can plug in a variety of databases.
- Allows standalone mode to work completely independently.

## Alerts

When inserting metrics into database, agent checks if the check previously
was successful or not. If the status of last check and current check changes,
it sends an alert.

## Page

An independent page can be created due to this architecture. Exporters
responsible for inserting metrics in the database can also fetch metrics
from the database. Agent offers a simple API endpoint through which
metrics can be fetched and a status page can be created.
