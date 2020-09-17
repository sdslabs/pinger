---
title: Architecture
nav_order: 2
has_children: true

---

# Architecture Overview

Pinger is made of four major components – client, app server, central
controller and agent. Before we dig deep into the architecture, we should become
familiar with the following keywords:

| Keyword                | What it is                                                                                                   |
|------------------------|--------------------------------------------------------------------------------------------------------------|
| **Check**              | Anything that can be pinged. _Say, pinging google.com using HTTP protocol and checking for status code 200._ |
| **Metric**             | Result of a ping. _Whether google.com returned status 200 in reasonable time or not._                        |
| **Controller**         | Something that runs a specific task again and again at regular intervals of time.                            |
| **Manager**            | Manages multiple controllers together and collects statistics from each of them.                             |
| **Page** (Status Page) | A collection of various checks and their corresponding metrics.                                              |

The basic idea: the client adds checks by requesting the app server. The app
server adds the checks to database. Central controller watches for changes
in the database and assigns the check to one of the agents registered
with the application.

![architecture-overview](/assets/images/architecture-overview.jpg)

Each component has more work than what can be seen from the diagram above.
To know more, let's dive into each component individually.
