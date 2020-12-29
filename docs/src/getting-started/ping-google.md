# Ping Google

Here, we'll deploy a status page using Pinger in the standalone mode, i.e.,
without complex deployment of the complete platform. But before that we need
to learn how to run checks.

## What is a Check?

Check, in hindsight, is a request sent to specified target which has its
response verified through some conditions which determine whether the request
returned valid response or not. These conditions might vary depending upon
the protocol of the request. For example, an HTTP request can be checked
from the status code received in the response.

## Writing the configuration file

For specifying the checks in the status page, we need to write a config file.
We will start by writing an ICMP check for **google.com**. This will be
equivalent of running the command:

```shell
$ ping google.com
```

The config file can be written in any of the following formats – YAML, TOML
or JSON. To create the status page we will be using the `agent` command which
takes the default file path to be `./agent.yml` so create a file with the
same name.

```yaml
# agent.yml

# We need to tell Pinger to run the agent in standalone mode since the
# default behaviour is something else.
standalone: true

# Configuration for where metrics are stored.
metrics:
  backend: log # We will just log them on the console for now

# Interval after which metrics are logged into database.
interval: 5s

# All the checks we need to run.
checks:
  - id: ping-google # unique ID
    name: Ping Google # human-readable name
    interval: 3s # Ping every 3 seconds
    timeout: 0.5s # Timeout if it takes longer than half a second
    input:
      type: ICMP # Protocol
    output:
      type: TIMEOUT # Condition for success/failure
    target: # Target to hit/request
      type: ADDRESS
      value: google.com
```

## Running the check

Now that we have the config file ready, we can launch our agent to ping the
Google servers. Assuming you have the Pinger binary, run the following
command:

```shell
$ path/to/pinger agent
INFO[0005] metrics for check (ping-google) Ping Google   check_id=ping-google check_name="Ping Google" duration=79.065304ms is_successful=true is_timeout=false start_time="2020-12-29 22:28:09.313387 +0530 IST m=+0.024071573"
INFO[0005] metrics for check (ping-google) Ping Google   check_id=ping-google check_name="Ping Google" duration=43.521461ms is_successful=true is_timeout=false start_time="2020-12-29 22:28:12.314517 +0530 IST m=+3.025195953"
INFO[0010] metrics for check (ping-google) Ping Google   check_id=ping-google check_name="Ping Google" duration=49.180263ms is_successful=true is_timeout=false start_time="2020-12-29 22:28:15.318509 +0530 IST m=+6.029181599"
INFO[0010] metrics for check (ping-google) Ping Google   check_id=ping-google check_name="Ping Google" duration=500ms is_successful=false is_timeout=true start_time="2020-12-29 22:28:18.318521 +0530 IST m=+9.029187048"
```

You will get an output similar to the one above. Each 5 seconds we get logs
for the metrics that we collected. We can also see the start times for each
log is 3 seconds apart, which is what we set for our check. Finally, the last
check is considered as a failure, given it took more than the set limit of
half a second (or 500 milli-second).
