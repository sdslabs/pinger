# Storing Metrics

Currently we only logged the metrics on the console. In long term, we would
want to persist the metrics in a time-series database. Let's see how.

## Setting up TimescaleDB

Pinger supports [TimescaleDB](https://www.timescale.com/) as of now.
Timescale is a PostgreSQL extension for storing time-series data.

Using the instructions
[here](https://docs.timescale.com/latest/getting-started/installation)
we can set it up. For this tutorial we can use Docker to spawn a container
which is Timescale ready.

```sh
$ docker run -d --name timescaledb -p 5432:5432 \
  -e POSTGRES_PASSWORD=password timescale/timescaledb:2.0.0-pg12
```

This will start a PostgreSQL instance on `:5432` with Timescale installed,
the user `postgres` having password `password`.

Before we configure our storage, we should create a database to store the
data in.

```sh
# This step is just to interactively exec into the container.
# If Timescale was set up natively into the system, it is not required.
$ docker exec -it 7e13ffbb3612 /bin/bash
# Open the postgres shell
$ psql -U postgres postgres
# Create a database named `pinger`
$ CREATE DATABASE pinger;
# List databases and verify that it is created
$ \l
                                   List of databases
      Name      |  Owner   | Encoding |  Collate   |   Ctype    |   Access privileges
----------------+----------+----------+------------+------------+-----------------------
 pinger         | postgres | UTF8     | en_US.utf8 | en_US.utf8 |
 postgres       | postgres | UTF8     | en_US.utf8 | en_US.utf8 |
 template0      | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
                |          |          |            |            | postgres=CTc/postgres
 template1      | postgres | UTF8     | en_US.utf8 | en_US.utf8 | =c/postgres          +
                |          |          |            |            | postgres=CTc/postgres
(4 rows)

```

Now that we have created our database, we can configure the agent to use it.

## Configuring storage backend

Now, we can update our config file replacing the log metrics backend with
the timescale instance.

```yaml
# agent.yml

# ...

# The new metrics configuration should look like this
metrics:
  backend: timescale
  host: 127.0.0.1
  port: 5432
  username: postgres
  password: password
  db_name: pinger
  ssl_mode: false # Let's just keep it off for now

# ...
```

That's it. We can restart our agent and see metrics being stored into the
database.

```sh
$ path/to/pinger agent
```

Now we can change the database from Postgres shell and see if metrics were
actually collected into the database.

```sh
$ \c pinger
You are now connected to database "pinger" as user "postgres".
$ SELECT * FROM metrics;
  check_id   | check_name  |          start_time           | duration  | timeout | success
-------------+-------------+-------------------------------+-----------+---------+---------
 ping-google | Ping Google | 2020-12-29 19:11:06.926541+00 |  72592580 | f       | t
 ping-google | Ping Google | 2020-12-29 19:11:09.931472+00 | 116203210 | f       | t
 ping-google | Ping Google | 2020-12-29 19:11:12.930745+00 |  54831874 | f       | t
 ping-google | Ping Google | 2020-12-29 19:11:15.930645+00 |  39025993 | f       | t
(4 rows)
```

You should see something like above.

Hurray! We have successfully set up a persistent storage backend for our
metrics.
