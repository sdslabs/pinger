---
title: Central Controller
parent: Architecture
nav_order: 3

---

# Central controller

Central component is present for assigning checks to agents. Since there
can be multiple agents across many machines, there is a need for a central
entity to resolve which agent should get which check.

## How it works

Central controller watches for changes in the database. If a check is created
or deleted or updated, it assigns or removes the check from the corresponding
agent. The checks are assigned to the agent with the minimum number of checks
at a given instance of time. Since it's rare for state of a check to update
often, this approach works.

Agents are registered with the database, which central controller watches.
Similar to checks, it manages agents. If an agent is deleted, it's checks are
reassigned to another agent. At any point there has to be atleast one agent.

The controller caches state of checks and agents in memory for quick access.
